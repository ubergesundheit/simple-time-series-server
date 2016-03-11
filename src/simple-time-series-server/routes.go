package main

import (
	"fmt"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
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
