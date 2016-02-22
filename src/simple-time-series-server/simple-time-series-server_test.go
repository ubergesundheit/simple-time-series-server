package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestValidateEntryEmpty(t *testing.T) {
	_, err := ValidateAndConvertEntry(Entry{})
	if err == nil {
		t.Fatal("should fail, empty Entry")
	}
}

func TestValidateEntryMissingFields(t *testing.T) {
	_, err := ValidateAndConvertEntry(Entry{Timestamp: time.Now(), Collection: "test_coll"})
	if err == nil {
		t.Fatal("should fail")
	}
}

func TestValidateEntryEmptyData(t *testing.T) {
	_, err := ValidateAndConvertEntry(Entry{Timestamp: time.Now(), Collection: "test_coll", Data: make(map[string]interface{})})
	if err == nil {
		t.Fatal("should fail")
	}
}

func TestInitDB(t *testing.T) {
	testFileName := "./test-initdb_simple-time-series-db.sqlite"
	app := App{}
	app.InitDB(testFileName)
	defer os.Remove(testFileName)

	var name string
	err := app.DB.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table';").Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsertDB(t *testing.T) {
	testFileName := "./test-initdb_simple-time-series-db.sqlite"
	app := App{}
	app.InitDB(testFileName)
	defer os.Remove(testFileName)

	var testjson = []byte(`{"collection":"testcollection","timestamp":"2016-01-07T23:57:10Z","data": {"foo": "bar", "f00": "baz"}}`)

	var entry Entry
	err := json.Unmarshal(testjson, &entry)
	if err != nil {
		t.Error(err)
	}
	err = app.CreateEntryInDB(entry)
	if err != nil {
		t.Error(err)
	}
	var collection string
	var timestamp int64
	var data []byte
	err = app.DB.QueryRow("SELECT collection, timestamp, data FROM entries;").Scan(&collection, &timestamp, &data)
	if err != nil {
		t.Error(err)
	}
}
