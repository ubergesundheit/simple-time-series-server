package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	DB *sql.DB
}

func main() {
	// flag handling
	handleFlags()

	app := App{}
	app.StartServer(address, filename)
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

	_, err = app.DB.Exec(CREATE_TABLE)
	checkErr(err)

	_, err = app.DB.Exec(CREATE_INDEX_COLLECTION)
	checkErr(err)

	if safeMode == false {
		_, err = app.DB.Exec(PRAGMAS)
		checkErr(err)
	}
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
