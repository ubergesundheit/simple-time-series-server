package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/boltdb/bolt"
)

type App struct {
	DB *bolt.DB
}

func main() {
	// flag handling
	handleConfig()

	app := App{}

	if flag.Arg(0) != "" && flag.Arg(1) != "" && flag.Arg(0) == "import" {
		app.Import(Filename, flag.Arg(1))
		os.Exit(0)
	}

	fmt.Println("Starting Server at", Address)
	app.StartServer(Address, Filename)
	defer app.DB.Close()
}

func (app *App) StartServer(addr string, dbFileName string) {
	err := app.InitDB(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
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
		rest.Get("/last", app.GetLast),
		rest.Get("/all", app.GetAll),
		rest.Post("/postEntry", app.PostEntry),
		//rest.Post("/postSignedEntry", app.PostJwtEntry),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(addr, api.MakeHandler()))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
