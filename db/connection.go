package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func GetDBConnection() *sql.DB {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	var annotationDbPath string = fmt.Sprintf("file:%s/Library/Containers/com.apple.iBooksX/Data/Documents/AEAnnotation/AEAnnotation_v10312011_1727_local.sqlite?cache=shared&mode=ro", homedir)
	var bookDbPath string = fmt.Sprintf("file:%s/Library/Containers/com.apple.iBooksX/Data/Documents/BKLibrary/BKLibrary-1-091020131601.sqlite?cache=shared&mode=ro", homedir)

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s", bookDbPath))
	if err != nil {
		log.Fatal(err)
	}

	// Attach second SQLLite database file to connection
	_, err = db.Exec(fmt.Sprintf("attach database '%s' as a", annotationDbPath))
	if err != nil {
		log.Println(fmt.Sprintf("attach database '%s' as a", annotationDbPath))
		log.Fatal(err)
	}

	return db
}
