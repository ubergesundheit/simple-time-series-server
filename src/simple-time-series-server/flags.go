package main

import (
	"flag"
)

const (
	filenameDefault       = "./simple-time-series-db.sqlite"
	filenameUsage         = "filename of the sqlite database"
	addressDefault        = "127.0.0.1:8080"
	addressUsage          = "adress to bind to"
	sqliteSafeModeDefault = true
	sqliteSafeModeUsage   = "If disabled, applies PRAGMA locking_mode = EXCLUSIVE;PRAGMA synchronous = OFF;PRAGMA journal_mode = OFF"
)

var filename string
var address string
var safeMode bool

func handleFlags() {
	flag.StringVar(&filename, "filename", filenameDefault, filenameUsage)
	flag.StringVar(&filename, "f", filenameDefault, filenameUsage)

	flag.StringVar(&address, "address", addressDefault, addressUsage)
	flag.StringVar(&address, "A", addressDefault, addressUsage)

	flag.BoolVar(&safeMode, "sqlite-safe-mode", sqliteSafeModeDefault, sqliteSafeModeUsage)
	flag.BoolVar(&safeMode, "s", sqliteSafeModeDefault, sqliteSafeModeUsage)

	flag.Parse()
}
