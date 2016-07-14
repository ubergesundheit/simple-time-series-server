package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func (app *App) Import(dbFilename string, importFilename string) {
	if DropDB == true {
		fmt.Println("Removing db file`", dbFilename, "`")
		os.Remove(dbFilename)
	} else {
		fmt.Println("Going to append to db `", dbFilename, "`")
	}

	fmt.Println("reading file `", importFilename, "` for import")
	file, e := ioutil.ReadFile(importFilename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}

	var entries []Entry
	json.Unmarshal(file, &entries)

	fmt.Println("found", len(entries), "entries, begin import")

	app.InitDB(dbFilename)

	tx, err := app.DB.Begin(true)
	checkErr(err)
	defer tx.Rollback()

	for _, entry := range entries {
		insertableEntry, err := ValidateAndConvertEntry(entry)
		checkErr(err)
		b, err := tx.CreateBucketIfNotExists([]byte(insertableEntry.Collection))
		checkErr(err)
		b.Put(insertableEntry.Timestamp, insertableEntry.Data)
		checkErr(err)
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		panic(err)
	}

	fmt.Println("imported", len(entries), "entries from `", importFilename, "` to `", dbFilename, "`")

	app.DB.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
