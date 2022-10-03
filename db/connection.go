package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func GetDBConnection() *sql.DB {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	annotationDbSearchPatch := fmt.Sprintf("%s/Library/Containers/com.apple.iBooksX/Data/Documents/AEAnnotation", homedir)
	booksDbSearchPatch := fmt.Sprintf("%s/Library/Containers/com.apple.iBooksX/Data/Documents/BKLibrary", homedir)
	annotationsFname := findByExt(annotationDbSearchPatch)
	booksFname := findByExt(booksDbSearchPatch)

	var annotationDbPathWithoutPrefix string = fmt.Sprintf("%s/%s", annotationDbSearchPatch, annotationsFname)
	var bookDbPathWithoutPrefix string = fmt.Sprintf("%s/%s", booksDbSearchPatch, booksFname)

	var annotationDbPathWithPrefix string = fmt.Sprintf("file:%s/%s", annotationDbSearchPatch, annotationsFname)
	var bookDbPathWithPrefix string = fmt.Sprintf("file:%s/%s", booksDbSearchPatch, booksFname)

	if _, err := os.Stat(annotationDbPathWithoutPrefix); errors.Is(err, os.ErrNotExist) {
		log.Fatal("iBooks files are not found.")
	}
	if _, err := os.Stat(bookDbPathWithoutPrefix); errors.Is(err, os.ErrNotExist) {
		log.Fatal("iBooks files are not found.")
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s", bookDbPathWithPrefix))
	if err != nil {
		log.Fatal(err)
	}

	// Attach second SQLLite database file to connection
	_, err = db.Exec(fmt.Sprintf("attach database '%s' as a", annotationDbPathWithPrefix))
	if err != nil {
		log.Println(fmt.Sprintf("attach database '%s' as a", annotationDbPathWithPrefix))
		log.Fatal(err)
	}

	return db
}

func findByExt(path string) string {
	ext := ".sqlite$"
	var fname string
	filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				fname = f.Name()
			}
		}
		return nil
	})

	return fname
}
