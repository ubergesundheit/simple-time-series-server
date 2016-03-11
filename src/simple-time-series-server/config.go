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

var Filename string
var Address string
var SafeMode bool

func handleConfig() {
	flag.StringVar(&Filename, "filename", filenameDefault, filenameUsage)
	flag.StringVar(&Filename, "f", filenameDefault, filenameUsage)

	flag.StringVar(&Address, "address", addressDefault, addressUsage)
	flag.StringVar(&Address, "A", addressDefault, addressUsage)

	flag.BoolVar(&SafeMode, "sqlite-safe-mode", sqliteSafeModeDefault, sqliteSafeModeUsage)
	flag.BoolVar(&SafeMode, "s", sqliteSafeModeDefault, sqliteSafeModeUsage)

	flag.Parse()
}
