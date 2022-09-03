package main

import (
	"database/sql"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const getAllBooksDbQueryConstant = "select ZBKLIBRARYASSET.ZASSETID, ZBKLIBRARYASSET.ZTITLE," +
	" count(annotations.ZAEANNOTATION.Z_PK" +
	" ZBKLIBRARYASSET.ZAUTHOR from ZBKLIBRARYASSET left join annotations.ZAEANNOTATION" +
	"annotations.ZAEANNOTATION.ZANNOTATIONASSETID = ZBKLIBRARYASSET.ZASSETID " +
	"GROUP BY ZBKLIBRARYASSET.ZASSETID"

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var annotationDbPath string = fmt.Sprintf("%s/Library/Containers/com.apple.iBooksX/Data/Documents/AEAnnotation/AEAnnotation_v10312011_1727_local.sqlite", homedir)
	var bookDbPath string = fmt.Sprintf("%s/Library/Containers/com.apple.iBooksX/Data/Documents/BKLibrary/BKLibrary-1-091020131601.sqlite", homedir)

	fmt.Println(annotationDbPath)
	fmt.Println(bookDbPath)

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s", bookDbPath))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	db.SetMaxOpenConns(1)
	//fmt.Println("Gamarjoba")

	// Attach second SQLLite database file to connection
	db.Exec(fmt.Sprintf("attach database %s as annotations", annotationDbPath))

	// Getting a list of books
	var (
		book_id     string
		book_title  string
		book_author string
	)
	rows, err := db.Query(getAllBooksDbQueryConstant)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Render table with books
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Book ID", "Title and Author"})

	for rows.Next() {
		err := rows.Scan(&book_id, &book_title, &book_author)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(book_id, book_title, book_author)
		t.AppendRows([]table.Row{
			//{1, "Arya", "Stark", 3000},
			{book_id, fmt.Sprintf("%s — %s", book_title, book_author)},
		})

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	t.Render()

}
