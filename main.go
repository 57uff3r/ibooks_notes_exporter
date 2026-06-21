package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v3"
	dbThings "ibooks_notes_exporter/db"
	_ "modernc.org/sqlite"
)

func main() {
	cmd := &cli.Command{
		Name:    "ibooks_notes_exporter",
		Usage:   "Export your records from Apple iBooks",
		Version: "v0.0.7",
		Commands: []*cli.Command{
			{
				Name:   "books",
				Usage:  "Get list of the books with notes and highlights",
				Action: getListOfBooks,
			},
			{
				Name:  "version",
				Usage: "Print version",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Printf("%s\n", cmd.Root().Version)
					return nil
				},
			},
			{
				Name:  "mcp",
				Usage: "Start MCP server for AI agents (stdio transport)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return runMCPServer()
				},
			},
			{
				Name:      "export",
				Usage:     "Export all notes and highlights from book with [BOOK_ID]",
				UsageText: "ibooks_notes_exporter export --book_id BOOK_ID_GOES_HERE",
				Action:    exportNotesAndHighlights,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "book_id",
						Required: true,
					},
					&cli.IntFlag{
						Name:  "skip_first_x_notes",
						Value: 0,
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func GetLastName(name string) string {
	words := strings.Fields(name)

	var lastName string
	for i := len(words) - 1; i >= 0; i-- {
		if !isHonorific(words[i]) {
			lastName = words[i]
			break
		}
	}

	lastName = strings.TrimSuffix(lastName, ",")
	lastName = strings.TrimSuffix(lastName, ".")

	return "(" + lastName + ")"
}

func isHonorific(word string) bool {
	return len(word) <= 3 && unicode.IsUpper(rune(word[0])) && (word[len(word)-1] == '.' || word[len(word)-1] == ',')
}

func GetLastNames(names string) string {
	nameList := strings.Split(names, " & ")

	if len(nameList) == 1 {
		return GetLastName(nameList[0])
	}

	if len(nameList) == 2 {
		return GetLastName(nameList[0]) + " & " + GetLastName(nameList[1])
	}

	firstName := nameList[0]
	lastNames := make([]string, len(nameList)-1)
	for i, name := range nameList[1:] {
		lastNames[i] = GetLastName(name)
	}
	return GetLastName(firstName) + " & " + strings.Join(lastNames, " & ")
}

func getListOfBooks(ctx context.Context, cmd *cli.Command) error {
	db := dbThings.GetDBConnection()
	defer db.Close()

	rows, err := db.Query(dbThings.GetAllBooksDbQueryConstant)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"SingleBook ID", "# notes", "Title and Author"})

	var singleBook dbThings.SingleBookInList
	for rows.Next() {
		err := rows.Scan(&singleBook.Id, &singleBook.Title, &singleBook.Author, &singleBook.Number)
		if err != nil {
			log.Fatal(err)
		}
		truncatedTitle := singleBook.Title
		if len(singleBook.Title) > 30 {
			truncatedTitle = singleBook.Title[:30] + "..."
		}
		standardizedAuthor := GetLastNames(singleBook.Author)
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

func exportNotesAndHighlights(ctx context.Context, cmd *cli.Command) error {
	db := dbThings.GetDBConnection()
	defer db.Close()

	bookId := cmd.String("book_id")
	skipXNotes := cmd.Int("skip_first_x_notes")
	fmt.Println(bookId)

	var book dbThings.SingleBook
	row := db.QueryRow(dbThings.GetBookDataById, bookId)
	err := row.Scan(&book.Name, &book.Author)
	if err != nil {
		log.Println(err)
		log.Fatal("SingleBook is not found in iBooks!")
	}

	fmt.Println(fmt.Sprintf("# %s — %s\n", book.Name, book.Author))

	rows, err := db.Query(dbThings.GetNotesHighlightsById, bookId, skipXNotes)
	if err != nil {
		log.Fatal(err)
	}

	var singleHightLightNote dbThings.SingleHighlightNote
	for rows.Next() {
		err := rows.Scan(&singleHightLightNote.HightLight, &singleHightLightNote.Note, &singleHightLightNote.Style)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(fmt.Sprintf("> <span style='background-color:%s;color:black'>%s</span>",styleToColor(singleHightLightNote.Style),strings.Replace(singleHightLightNote.HightLight, "\n", "", -1)))

		if singleHightLightNote.Note.Valid {
			fmt.Println(fmt.Sprintf("\n%s", strings.Replace(singleHightLightNote.Note.String, "\n", "", -1)))
		}

		fmt.Print("---\n\n\n")
	}

	return nil
}

func styleToColor(style int) string {
	switch style {
	case 1: 
		// green
		return "#a8e196"
	case 2:
		// blue
		return "#a5c3ff" 
	case 3: 
		// yellow
		return "#fde15c"
	case 4: 
		// pink
		return "#ffaabf"	
	case 5: 
		// purple
		return "#cdbbfb"	
	default:
		return ""		
	}
}