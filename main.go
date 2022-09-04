package main

import (
	"database/sql"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const getAllBooksDbQueryConstant = `
	select 
		ZBKLIBRARYASSET.ZASSETID,
		ZBKLIBRARYASSET.ZTITLE,
		ZBKLIBRARYASSET.ZAUTHOR,    
		count(a.ZAEANNOTATION.Z_PK)
	from ZBKLIBRARYASSET left join a.ZAEANNOTATION
		on a.ZAEANNOTATION.ZANNOTATIONASSETID = ZBKLIBRARYASSET.ZASSETID
	WHERE a.ZAEANNOTATION.ZANNOTATIONSELECTEDTEXT NOT NULL
	GROUP BY ZBKLIBRARYASSET.ZASSETID;
`

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var annotationDbPath string = fmt.Sprintf("file:%s/Library/Containers/com.apple.iBooksX/Data/Documents/AEAnnotation/AEAnnotation_v10312011_1727_local.sqlite?cache=shared&mode=ro", homedir)
	var bookDbPath string = fmt.Sprintf("file:%s/Library/Containers/com.apple.iBooksX/Data/Documents/BKLibrary/BKLibrary-1-091020131601.sqlite?cache=shared&mode=ro", homedir)

	fmt.Println(annotationDbPath)
	fmt.Println(bookDbPath)

	//sql.Register("sqlite3_hooked",
	//	&sqlite3.SQLiteDriver{
	//		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
	//			conn.Exec("ATTACH DATABASE '"+annotationDbPath+"' AS 'annotations';", nil)
	//			return nil
	//		},
	//	})

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s", bookDbPath))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Attach second SQLLite database file to connection
	_, err = db.Exec(fmt.Sprintf("attach database '%s' as a", annotationDbPath))
	if err != nil {
		log.Println(fmt.Sprintf("attach database '%s' as a", annotationDbPath))
		log.Fatal(err)
	}

	// Getting a list of books
	var (
		book_id     string
		book_title  string
		book_author string
		number      int
	)
	rows, err := db.Query(getAllBooksDbQueryConstant)
	if err != nil {
		log.Println(getAllBooksDbQueryConstant)
		log.Fatal(err)
	}
	defer rows.Close()

	// Render table with books
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Book ID", "Number of notes", "Title and Author"})

	for rows.Next() {
		err := rows.Scan(&book_id, &book_title, &book_author, &number)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(book_id, book_title, book_author)
		t.AppendRows([]table.Row{
			//{1, "Arya", "Stark", 3000},
			{book_id, number, fmt.Sprintf("%s — %s", book_title, book_author)},
		})

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	t.Render()

}
