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

	annotationDbSearchPath := fmt.Sprintf("%s/Library/Containers/com.apple.iBooksX/Data/Documents/AEAnnotation", homedir)
	booksDbSearchPath := fmt.Sprintf("%s/Library/Containers/com.apple.iBooksX/Data/Documents/BKLibrary", homedir)
	annotationsFname := findByExt(annotationDbSearchPath)
	booksFname := findByExt(booksDbSearchPath)

	annotationDbPath := fmt.Sprintf("%s/%s", annotationDbSearchPath, annotationsFname)
	bookDbPath := fmt.Sprintf("%s/%s", booksDbSearchPath, booksFname)

	if _, err := os.Stat(annotationDbPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("iBooks files are not found.")
	}
	if _, err := os.Stat(bookDbPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("iBooks files are not found.")
	}

	db, err := sql.Open("sqlite", bookDbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(fmt.Sprintf("attach database '%s' as a", annotationDbPath))
	if err != nil {
		log.Println(fmt.Sprintf("attach database '%s' as a", annotationDbPath))
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
