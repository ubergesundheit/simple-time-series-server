package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
)

func (app *App) Import(dbFilename string, importFilename string) {
	fmt.Println("reading file `", importFilename, "` for import")
	file, e := ioutil.ReadFile(importFilename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}

	var entries []Entry
	json.Unmarshal(file, &entries)

	fmt.Println("found", len(entries), "entries, begin import")

	app.InitDB(dbFilename)

	_, err := app.DB.Exec("BEGIN TRANSACTION;")
	checkErr(err)

	for _, entry := range entries {
		app.CreateEntryInDB(entry)
	}
	fmt.Println("imported", len(entries), "entries from `", importFilename, "` to `", dbFilename, "`")

	_, err = app.DB.Exec("COMMIT TRANSACTION;")
	checkErr(err)

	app.DB.Close()
}
