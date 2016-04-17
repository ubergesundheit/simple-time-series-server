package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	BodyIsEmptyError = "Body is empty"
)

// this route is some kind of an allrounder
// it parses the following parameters:
// collections -> restrict collections, separated by comma, if omitted, include all Collections
// from -> which timespan to include, parsed by time.ParseDuration, if omitted -> -24h
func (app *App) GetLast(w rest.ResponseWriter, r *rest.Request) {
	parameters := r.URL.Query()
	from := time.Duration(-24) * time.Hour
	collections := ""

	if parameters["from"] != nil {
		duration, err := time.ParseDuration(parameters["from"][0])
		if err != nil {
			fmt.Println(err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		from = duration
	}

	if parameters["collections"] != nil {
		collections = parameters["collections"][0]
	}

	response, err := app.GetLastFromDB(collections, from)
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(response) != 0 {
		w.WriteJson(response)
	} else {
		w.WriteJson([]string{})
	}
}

func (app *App) GetAll(w rest.ResponseWriter, r *rest.Request) {
	parameters := r.URL.Query()
	collections := ""

	if parameters["collections"] != nil {
		collections = parameters["collections"][0]
	}

	response, err := app.GetAllFromDB(collections)
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(response) != 0 {
		w.WriteJson(response)
	} else {
		w.WriteJson([]string{})
	}
}

func (app *App) GetLatest(w rest.ResponseWriter, r *rest.Request) {
	response, err := app.GetLatestFromDB()
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(response) != 0 {
		w.WriteJson(response)
	} else {
		w.WriteJson([]string{})
	}

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

func (app *App) PostJwtEntry(w rest.ResponseWriter, r *rest.Request) {
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if len(content) == 0 {
		fmt.Println(err.Error())
		rest.Error(w, BodyIsEmptyError, http.StatusInternalServerError)
	}

	parsedEntry, err := ParseJwt(string(content))
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = app.CreateEntryInDB(parsedEntry)
	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
