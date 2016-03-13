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
	testFileName := "./test-initdb_simple-time-series-db.db"
	app := App{}
	err := app.InitDB(testFileName)
	defer os.Remove(testFileName)

	if err != nil {
		t.Fatal(err)
	}
}

func TestInsertDB(t *testing.T) {
	testFileName := "./test-insertdb_simple-time-series-db.db"
	app := App{}
	err := app.InitDB(testFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFileName)

	var testjson = []byte(`{"collection":"testcollection","timestamp":"2016-01-07T23:57:10Z","data": {"foo": "bar", "f00": "baz"}}`)

	var entry Entry
	err = json.Unmarshal(testjson, &entry)
	if err != nil {
		t.Error(err)
	}
	err = app.CreateEntryInDB(entry)
	if err != nil {
		t.Error(err)
	}
}
