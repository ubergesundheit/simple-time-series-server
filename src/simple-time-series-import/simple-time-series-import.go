package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	Collection string      `json:"collection"`
	Timestamp  time.Time   `json:"timestamp"`
	Data       interface{} `json:"data"`
}

func main() {
	file, e := ioutil.ReadFile("./all")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
	}

	var entries []Entry
	json.Unmarshal(file, &entries)

	db, err := sql.Open("sqlite3", "./simple-time-series-db.sqlite")
	checkErr(err)

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS entries (collection TEXT, timestamp INTEGER, data BLOB);")
	checkErr(err)

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS index_collection ON entries (collection);")
	checkErr(err)

	_, err = db.Exec("PRAGMA locking_mode = EXCLUSIVE;PRAGMA synchronous = OFF;PRAGMA journal_mode = OFF;")
	checkErr(err)

	_, err = db.Exec("BEGIN TRANSACTION;")
	checkErr(err)

	for _, entry := range entries {
		jsonbytes, err := json.Marshal(entry.Data)
		checkErr(err)

		stmt, err := db.Prepare("INSERT INTO `entries` (collection, timestamp, data) VALUES (?, ?, ?);")
		checkErr(err)

		_, err = stmt.Exec(entry.Collection, entry.Timestamp.UTC().Unix(), jsonbytes)
		checkErr(err)
	}

	_, err = db.Exec("COMMIT TRANSACTION;")
	checkErr(err)

	db.Close()
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
