package main

import (
	"encoding/json"
	"errors"
	"time"
)

type Entry struct {
	Collection string                 `json:"collection"`
	Timestamp  time.Time              `json:"timestamp"`
	Data       map[string]interface{} `json:"data"`
}

type InsertableEntry struct {
	Collection string
	Timestamp  int64
	Data       []byte
}

func ValidateAndConvertEntry(e Entry) (InsertableEntry, error) {
	if len(e.Collection) == 0 || e.Timestamp.IsZero() == true || len(e.Data) == 0 {
		return InsertableEntry{}, errors.New("fields `collection`, `timestamp` and `data` are required and must be non-empty")
	}

	jsonData, err := json.Marshal(e.Data)
	if err != nil {
		return InsertableEntry{}, err
	}

	ie := InsertableEntry{
		Collection: e.Collection,
		Timestamp:  e.Timestamp.UTC().Unix(),
		Data:       jsonData,
	}

	return ie, nil
}
