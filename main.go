package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	dbThings "ibooks_notes_exporter/db"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "Ibooks notes exporter",
		Usage:   "Export your records from Apple iBooks",
		Authors: []*cli.Author{{Name: "Andrey Korchak", Email: "me@akorchak.software"}},
		Commands: []*cli.Command{
			{
				Name:   "books",
				Usage:  "Get list of the books with notes and highlights",
				Action: getListOfBooks,
			},
			{
				Name:      "export",
				HideHelp:  false,
				Usage:     "Export all notes and highlights from book with [BOOK_ID]",
				UsageText: "Export all notes and highlights from book with [BOOK_ID]",
				Action:    exportNotesAndHighlights,
				ArgsUsage: "ibooks_notes_exporter export BOOK_ID_GOES_HERE",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "book_id"},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func getListOfBooks(cCtx *cli.Context) error {
	db := dbThings.GetDBConnection()

	// Getting a list of books
	rows, err := db.Query(dbThings.GetAllBooksDbQueryConstant)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Render table with books
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"SingleBook ID", "Number of notes", "Title and Author"})

	var singleBook dbThings.SingleBookInList
	for rows.Next() {
		err := rows.Scan(&singleBook.Id, &singleBook.Title, &singleBook.Author, &singleBook.Number)
		if err != nil {
			log.Fatal(err)
		}
		t.AppendRows([]table.Row{
			{singleBook.Id, singleBook.Number, fmt.Sprintf("%s — %s", singleBook.Title, singleBook.Author)},
		})
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	t.Render()
	return nil
}

func exportNotesAndHighlights(cCtx *cli.Context) error {
	db := dbThings.GetDBConnection()
	defer db.Close()

	if cCtx.Args().Len() != 1 {
		log.Fatal("For exporting notes and highlights, you have to pass BOOK_ID: ibooks_notes_exporter export BOOK_ID_GOES_HERE")
	}

	fmt.Println(cCtx.Args().Get(0))

	var book dbThings.SingleBook
	row := db.QueryRow(dbThings.GetBookDataById, cCtx.Args().Get(0))
	err := row.Scan(&book.Name, &book.Author)
	if err != nil {
		//log.Fatal()
		log.Println(err)
		log.Fatal("SingleBook is not found in iBooks!")
	}

	// Render MarkDown into STDOUT
	fmt.Println(fmt.Sprintf("# %s — %s\n", book.Name, book.Author))

	rows, err := db.Query(dbThings.GetNotesHighlightsById, cCtx.Args().Get(0))
	if err != nil {
		log.Fatal(err)
	}

	var singleHightLightNote dbThings.SingleHighlightNote
	for rows.Next() {
		err := rows.Scan(&singleHightLightNote.HightLight, &singleHightLightNote.Note)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(fmt.Sprintf("> %s", singleHightLightNote.HightLight))

		if singleHightLightNote.Note.Valid {
			fmt.Println(fmt.Sprintf("\n%s", singleHightLightNote.Note.String))
		}

		fmt.Println("---\n\n")

	}

	return nil
}
