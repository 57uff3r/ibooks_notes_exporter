package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	dbThings "ibooks_notes_exporter/db"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	app := &cli.App{
		Name:    "Ibooks notes exporter",
		Usage:   "Export your records from Apple iBooks",
		Authors: []*cli.Author{{Name: "Andrey Korchak", Email: "me@akorchak.software"}},
		Version: "v0.0.5",
		Commands: []*cli.Command{
			{
				Name:   "books",
				Usage:  "Get list of the books with notes and highlights",
				Action: getListOfBooks,
			},
			{
				Name: "version",
				Action: func(context *cli.Context) error {
					fmt.Printf("%s\n", context.App.Version)
					return nil
				},
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

func GetLastName(name string) string {
	// Split the input string into words
	words := strings.Fields(name)

	// Search backwards from the end of the string for the first non-title word
	var lastName string
	for i := len(words) - 1; i >= 0; i-- {
		if !isHonorific(words[i]) {
			lastName = words[i]
			break
		}
	}

	// Remove any trailing commas or periods from the last name
	lastName = strings.TrimSuffix(lastName, ",")
	lastName = strings.TrimSuffix(lastName, ".")

	// Return the last name in parentheses
	return "(" + lastName + ")"
}

// Helper function to check if a word is an honorific title
func isHonorific(word string) bool {
	return len(word) <= 3 && unicode.IsUpper(rune(word[0])) && (word[len(word)-1] == '.' || word[len(word)-1] == ',')
}

func GetLastNames(names string) string {
	// Split the input string into individual names
	nameList := strings.Split(names, " & ")

	// If there is only one name, just return the last name
	if len(nameList) == 1 {
		return GetLastName(nameList[0])
	}

	// If there are two names, combine the last names with "&"
	if len(nameList) == 2 {
		return GetLastName(nameList[0]) + " & " + GetLastName(nameList[1])
	}

	// If there are more than two names, combine the first name and last names with "&"
	firstName := nameList[0]
	lastNames := make([]string, len(nameList)-1)
	for i, name := range nameList[1:] {
		lastNames[i] = GetLastName(name)
	}
	return GetLastName(firstName) + " & " + strings.Join(lastNames, " & ")
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
	t.AppendHeader(table.Row{"SingleBook ID", "# notes", "Title and Author"})

	var singleBook dbThings.SingleBookInList
	for rows.Next() {
		err := rows.Scan(&singleBook.Id, &singleBook.Title, &singleBook.Author, &singleBook.Number)
		if err != nil {
			log.Fatal(err)
		}
		// truncate title as needed so that table doesn't wrap when terminal width is narrow
		truncatedTitle := singleBook.Title
		if len(singleBook.Title) > 30 {
			truncatedTitle = singleBook.Title[:30] + "..."
		}
		// shortened author name(s)
		standardizedAuthor := GetLastNames(singleBook.Author)
		// The title and author looks like: "My Great Book (Doe)"
		t.AppendRows([]table.Row{
			{singleBook.Id, singleBook.Number, fmt.Sprintf("%s %s", truncatedTitle, standardizedAuthor)},
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
