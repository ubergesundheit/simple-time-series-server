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

var SELECT_LATEST string = "SELECT `collection`, MAX(`timestamp`) as timestamp, `data` FROM `entries` GROUP BY collection;"
var INSERT_INTO string = "INSERT INTO `entries` (collection, timestamp, data) VALUES (?, ?, ?);"

type Impl struct {
	DB *sql.DB
}

type Entry struct {
	Collection string                 `json:"collection"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}

func main() {
	i := Impl{}
	i.InitDB()

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		AllowedMethods:        []string{"GET", "POST"},
		AllowedHeaders: []string{
			"Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           2592000,
	})
	router, err := rest.MakeRouter(
		rest.Get("/latest", i.GetLatest),
		rest.Post("/postEntry", i.CreateEntry),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
	defer i.DB.Close()
}

func (i *Impl) InitDB() {
	var err error
	i.DB, err = sql.Open("sqlite3", "./simple-time-series-db.sqlite")
	checkErr(err)

	_, err = i.DB.Exec("PRAGMA locking_mode = EXCLUSIVE;PRAGMA synchronous = OFF;PRAGMA journal_mode = OFF;")
	checkErr(err)
}

func (i *Impl) GetLatest(w rest.ResponseWriter, r *rest.Request) {
	rows, err := i.DB.Query(SELECT_LATEST)
	checkErr(err)
	var response []Entry

	for rows.Next() {
		var timestamp int64
		var collection string
		var blobdata []uint8
		var data map[string]interface{}
		err = rows.Scan(&collection, &timestamp, &blobdata)
		checkErr(err)
		err := json.Unmarshal(blobdata, &data)
		checkErr(err)
		response = append(response, Entry{
			Collection: collection,
			Timestamp:  time.Unix(timestamp, 0).UTC(),
			Data:       data,
		})
	}
	w.WriteJson(response)
}

func (i *Impl) CreateEntry(w rest.ResponseWriter, r *rest.Request) {
	newEntry := Entry{}
	err := r.DecodeJsonPayload(&newEntry)
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if newEntry.Collection == "" {
		rest.Error(w, "collection required", 400)
		return
	}
	if newEntry.Data == nil {
		rest.Error(w, "data required", 400)
		return
	}

	if len(newEntry.Data) == 0 {
		rest.Error(w, "data 1!", 400)
		return

	}

	jsonData, err := json.Marshal(newEntry.Data)
	checkErr(err)

	stmt, err := i.DB.Prepare(INSERT_INTO)
	checkErr(err)
	res, err := stmt.Exec(newEntry.Collection, newEntry.Timestamp.UTC().Unix(), jsonData)
	checkErr(err)

	_, err = res.LastInsertId()
	checkErr(err)
	w.WriteHeader(http.StatusCreated)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
