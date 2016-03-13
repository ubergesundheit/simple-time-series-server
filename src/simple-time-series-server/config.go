package main

import (
	"flag"
)

const (
	filenameDefault = "./simple-time-series-db.db"
	filenameUsage   = "filename of the sqlite database"
	addressDefault  = "127.0.0.1:8080"
	addressUsage    = "adress to bind to"
)

var Filename string
var Address string
var SafeMode bool

func handleConfig() {
	flag.StringVar(&Filename, "filename", filenameDefault, filenameUsage)
	flag.StringVar(&Filename, "f", filenameDefault, filenameUsage)

	flag.StringVar(&Address, "address", addressDefault, addressUsage)
	flag.StringVar(&Address, "A", addressDefault, addressUsage)

	flag.Parse()
}
