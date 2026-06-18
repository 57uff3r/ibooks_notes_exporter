package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	dbThings "ibooks_notes_exporter/db"
	_ "modernc.org/sqlite"
)

func runMCPServer() error {
	s := server.NewMCPServer(
		"ibooks-notes",
		"0.0.7",
		server.WithToolCapabilities(true),
	)

	listBooksTool := mcp.NewTool("list_books",
		mcp.WithDescription("List all books in Apple iBooks that have highlights or notes. Returns book IDs, titles, authors, and note counts."),
		mcp.WithReadOnlyHintAnnotation(true),
	)

	getNotesTool := mcp.NewTool("get_notes",
		mcp.WithDescription("Get all highlights and notes from a specific book. Returns Markdown with quoted highlights and inline notes."),
		mcp.WithString("book_id",
			mcp.Required(),
			mcp.Description("Book asset ID from list_books"),
		),
		mcp.WithReadOnlyHintAnnotation(true),
	)

	searchNotesTool := mcp.NewTool("search_notes",
		mcp.WithDescription("Search across all book highlights and notes for a keyword or phrase. Case-insensitive. Returns matching highlights with their book titles."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search term to find in highlights and notes"),
		),
		mcp.WithReadOnlyHintAnnotation(true),
	)

	s.AddTool(listBooksTool, handleListBooks)
	s.AddTool(getNotesTool, handleGetNotes)
	s.AddTool(searchNotesTool, handleSearchNotes)

	return server.ServeStdio(s)
}

type bookJSON struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	NumNotes int    `json:"num_notes"`
}

func handleListBooks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	db := dbThings.GetDBConnection()
	defer db.Close()

	rows, err := db.Query(dbThings.GetAllBooksDbQueryConstant)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.NewTextContent(fmt.Sprintf("Error querying books: %v", err))},
			IsError: true,
		}, nil
	}
	defer rows.Close()

	var books []bookJSON
	for rows.Next() {
		var b dbThings.SingleBookInList
		if err := rows.Scan(&b.Id, &b.Title, &b.Author, &b.Number); err != nil {
			continue
		}
		books = append(books, bookJSON{
			ID:       b.Id,
			Title:    b.Title,
			Author:   b.Author,
			NumNotes: b.Number,
		})
	}

	data, _ := json.Marshal(books)
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.NewTextContent(string(data))},
	}, nil
}

func handleGetNotes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	bookID, ok := request.GetArguments()["book_id"].(string)
	if !ok || bookID == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.NewTextContent("book_id is required")},
			IsError: true,
		}, nil
	}

	db := dbThings.GetDBConnection()
	defer db.Close()

	var book dbThings.SingleBook
	row := db.QueryRow(dbThings.GetBookDataById, bookID)
	if err := row.Scan(&book.Name, &book.Author); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.NewTextContent("Book not found in iBooks")},
			IsError: true,
		}, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s — %s\n\n", book.Name, book.Author))

	rows, err := db.Query(dbThings.GetNotesHighlightsById, bookID, 0)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.NewTextContent(fmt.Sprintf("Error querying notes: %v", err))},
			IsError: true,
		}, nil
	}
	defer rows.Close()

	for rows.Next() {
		var note dbThings.SingleHighlightNote
		if err := rows.Scan(&note.HightLight, &note.Note); err != nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("> %s\n", strings.ReplaceAll(note.HightLight, "\n", "")))
		if note.Note.Valid {
			sb.WriteString(fmt.Sprintf("\n%s\n", strings.ReplaceAll(note.Note.String, "\n", "")))
		}
		sb.WriteString("---\n\n")
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.NewTextContent(sb.String())},
	}, nil
}

const searchNotesQuery = `
	SELECT
		ZBKLIBRARYASSET.ZTITLE,
		ZBKLIBRARYASSET.ZAUTHOR,
		a.ZAEANNOTATION.ZANNOTATIONSELECTEDTEXT,
		a.ZAEANNOTATION.ZANNOTATIONNOTE
	FROM a.ZAEANNOTATION
	JOIN ZBKLIBRARYASSET
		ON a.ZAEANNOTATION.ZANNOTATIONASSETID = ZBKLIBRARYASSET.ZASSETID
	WHERE a.ZAEANNOTATION.ZANNOTATIONSELECTEDTEXT NOT NULL
		AND (
			a.ZAEANNOTATION.ZANNOTATIONSELECTEDTEXT LIKE '%' || $1 || '%'
			OR a.ZAEANNOTATION.ZANNOTATIONNOTE LIKE '%' || $1 || '%'
		)
	ORDER BY ZBKLIBRARYASSET.ZTITLE, a.ZAEANNOTATION.ZPLLOCATIONRANGESTART ASC
`

type searchResult struct {
	title     string
	author    string
	highlight string
	note      sql.NullString
}

func deduplicateResults(results []searchResult) []searchResult {
	var deduped []searchResult
	for _, r := range results {
		isDupe := false
		for j, d := range deduped {
			if d.title != r.title {
				continue
			}
			if strings.HasPrefix(r.highlight, d.highlight) || strings.HasPrefix(d.highlight, r.highlight) {
				isDupe = true
				if len(r.highlight) > len(d.highlight) {
					deduped[j] = r
				}
				break
			}
		}
		if !isDupe {
			deduped = append(deduped, r)
		}
	}
	return deduped
}

func handleSearchNotes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.GetArguments()["query"].(string)
	if !ok || query == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.NewTextContent("query is required")},
			IsError: true,
		}, nil
	}

	db := dbThings.GetDBConnection()
	defer db.Close()

	rows, err := db.Query(searchNotesQuery, query)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.NewTextContent(fmt.Sprintf("Error searching: %v", err))},
			IsError: true,
		}, nil
	}
	defer rows.Close()

	var results []searchResult
	for rows.Next() {
		var r searchResult
		if err := rows.Scan(&r.title, &r.author, &r.highlight, &r.note); err != nil {
			continue
		}
		results = append(results, r)
	}

	deduped := deduplicateResults(results)

	if len(deduped) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{mcp.NewTextContent("No results found")},
		}, nil
	}

	var sb strings.Builder
	for _, r := range deduped {
		sb.WriteString(fmt.Sprintf("## %s — %s\n\n", r.title, r.author))
		sb.WriteString(fmt.Sprintf("> %s\n", strings.ReplaceAll(r.highlight, "\n", "")))
		if r.note.Valid {
			sb.WriteString(fmt.Sprintf("\n%s\n", strings.ReplaceAll(r.note.String, "\n", "")))
		}
		sb.WriteString("---\n\n")
	}

	header := fmt.Sprintf("Found %d results for \"%s\":\n\n", len(deduped), query)
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.NewTextContent(header + sb.String())},
	}, nil
}
