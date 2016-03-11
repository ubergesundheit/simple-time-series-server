package main

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (app *App) InitDB(dbFileName string) error {
	var err error
	app.DB, err = sql.Open("sqlite3", dbFileName)
	if err != nil {
		return err
	}

	_, err = app.DB.Exec(CREATE_TABLE)
	if err != nil {
		return err
	}

	_, err = app.DB.Exec(CREATE_INDEX_COLLECTION)
	if err != nil {
		return err
	}

	if SafeMode == false {
		_, err = app.DB.Exec(PRAGMAS)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) GetLatestFromDB() ([]Entry, error) {
	rows, err := app.DB.Query(SELECT_LATEST)
	if err != nil {
		return []Entry{}, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var timestamp int64
		var collection string
		var blobdata []uint8
		var data map[string]interface{}
		err = rows.Scan(&collection, &timestamp, &blobdata)
		if err != nil {
			return []Entry{}, err
		}
		err := json.Unmarshal(blobdata, &data)
		if err != nil {
			return []Entry{}, err
		}
		entries = append(entries, Entry{
			Collection: collection,
			Timestamp:  time.Unix(timestamp, 0).UTC(),
			Data:       data,
		})
	}
	if err = rows.Close(); err != nil {
		return []Entry{}, err
	}
	err = rows.Err()
	if err != nil {
		return []Entry{}, err
	}
	return entries, nil
}

func (app *App) CreateEntryInDB(entry Entry) error {
	insertableEntry, err := ValidateAndConvertEntry(entry)
	if err != nil {
		return err
	}

	stmt, err := app.DB.Prepare(INSERT_INTO)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(insertableEntry.Collection, insertableEntry.Timestamp, insertableEntry.Data)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
