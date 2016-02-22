package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/mattn/go-sqlite3"
)

var SQL_STATEMENTS = map[string]string{
	"SELECT_LATEST":           "SELECT `collection`, MAX(`timestamp`) as timestamp, `data` FROM `entries` GROUP BY collection;",
	"INSERT_INTO":             "INSERT INTO `entries` (collection, timestamp, data) VALUES (?, ?, ?);",
	"CREATE_TABLE":            "CREATE TABLE IF NOT EXISTS entries (collection TEXT, timestamp INTEGER, data BLOB);",
	"CREATE_INDEX_COLLECTION": "CREATE INDEX IF NOT EXISTS index_collection ON entries (collection);",
	"PRAGMAS":                 "PRAGMA locking_mode = EXCLUSIVE;PRAGMA synchronous = OFF;PRAGMA journal_mode = OFF;",
}

type App struct {
	DB *sql.DB
}

type Entry struct {
	Collection string                 `json:"collection"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}

type InsertableEntry struct {
	Collection string
	Timestamp  int64
	Data       []byte
}

func main() {
	app := App{}
	app.StartServer(":8080", "./simple-time-series-db.sqlite")
	defer app.DB.Close()
}

func (app *App) StartServer(addr string, dbFileName string) {
	app.InitDB(dbFileName)
	api := rest.NewApi()
	api.Use(rest.DefaultProdStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		AllowedMethods:        []string{"GET", "POST"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           2592000,
	})
	router, err := rest.MakeRouter(
		rest.Get("/latest", app.GetLatest),
		rest.Post("/postEntry", app.PostEntry),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(addr, api.MakeHandler()))
}

func (app *App) InitDB(dbFileName string) {
	var err error
	app.DB, err = sql.Open("sqlite3", dbFileName)
	checkErr(err)

	_, err = app.DB.Exec(SQL_STATEMENTS["CREATE_TABLE"])
	checkErr(err)

	_, err = app.DB.Exec(SQL_STATEMENTS["CREATE_INDEX_COLLECTION"])
	checkErr(err)

	_, err = app.DB.Exec(SQL_STATEMENTS["PRAGMAS"])
	checkErr(err)
}

func (app *App) GetLatestFromDB() ([]Entry, error) {
	rows, err := app.DB.Query(SQL_STATEMENTS["SELECT_LATEST"])
	if err != nil {
		return []Entry{}, err
	}

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
	return entries, nil
}

func (app *App) GetLatest(w rest.ResponseWriter, r *rest.Request) {
	response, err := app.GetLatestFromDB()
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(response)
}

func (app *App) PostEntry(w rest.ResponseWriter, r *rest.Request) {
	newEntry := Entry{}
	err := r.DecodeJsonPayload(&newEntry)
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.CreateEntryInDB(newEntry)
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) CreateEntryInDB(entry Entry) error {
	insertableEntry, err := ValidateAndConvertEntry(entry)
	if err != nil {
		return err
	}

	stmt, err := app.DB.Prepare(SQL_STATEMENTS["INSERT_INTO"])
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

func ValidateAndConvertEntry(e Entry) (InsertableEntry, error) {
	if len(e.Collection) == 0 || e.Timestamp.IsZero() == true || len(e.Data) == 0 {
		return InsertableEntry{}, errors.New("fields `collection`, `timestamp` and `data` are required and must be non-empty")
	}

	jsonData, err := json.Marshal(e.Data)
	if err != nil {
		return InsertableEntry{}, err
	}

	ie := InsertableEntry{
		Collection: e.Collection,
		Timestamp:  e.Timestamp.UTC().Unix(),
		Data:       jsonData,
	}

	return ie, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
