package main

import (
	"flag"
)

const (
	filenameDefault = "./simple-time-series-db.db"
	filenameUsage   = "filename of the sqlite database"
	addressDefault  = "127.0.0.1:8080"
	addressUsage    = "adress to bind to"
	dropDefault     = false
	dropUsage       = "only for import command: Drop the database before importing"
)

var Filename string
var Address string
var DropDB bool

func handleConfig() {
	flag.StringVar(&Filename, "filename", filenameDefault, filenameUsage)
	flag.StringVar(&Filename, "f", filenameDefault, filenameUsage)

	flag.StringVar(&Address, "address", addressDefault, addressUsage)
	flag.StringVar(&Address, "A", addressDefault, addressUsage)

	flag.BoolVar(&DropDB, "drop", dropDefault, dropUsage)

	flag.Parse()
}
