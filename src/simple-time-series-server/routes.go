package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

var (
	BodyIsEmptyError = "Body is empty"
)

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
