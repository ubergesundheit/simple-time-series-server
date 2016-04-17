package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

func (app *App) InitDB(dbFileName string) error {
	var err error
	app.DB, err = bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) GetAllFromDB(collections string) ([]Entry, error) {
	var entries []Entry
	err := app.DB.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			strName := string(name)
			// yeah, this is not very beautiful, but it seems to work
			// just use the bucket name as substring for a strings.Index search on the collections string..
			if len(collections) == 0 || strings.Index(collections, strName) != -1 {
				cursor := bucket.Cursor()

				for timestamp, rawdata := cursor.First(); timestamp != nil; timestamp, rawdata = cursor.Next() {
					var data map[string]interface{}
					err := json.Unmarshal(rawdata, &data)
					if err != nil {
						return err
					}
					parsedTime, err := time.Parse(time.RFC3339, string(timestamp))
					if err != nil {
						return err
					}
					entries = append(entries, Entry{
						Collection: strName,
						Timestamp:  parsedTime,
						Data:       data,
					})
				}
			}

			return nil
		})
	})
	if err != nil {
		return []Entry{}, err
	}

	return entries, nil
}

func (app *App) GetLastFromDB(collections string, from time.Duration) ([]Entry, error) {
	var entries []Entry
	err := app.DB.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			strName := string(name)
			// yeah, this is not very beautiful, but it seems to work
			// just use the bucket name as substring for a strings.Index search on the collections string..
			if len(collections) == 0 || strings.Index(collections, strName) != -1 {
				cursor := bucket.Cursor()

				// compute the timestamp `from` and store in a
				startTimestamp := []byte(time.Now().UTC().Add(from).Format(time.RFC3339))

				for timestamp, rawdata := cursor.Seek(startTimestamp); timestamp != nil; timestamp, rawdata = cursor.Next() {
					var data map[string]interface{}
					err := json.Unmarshal(rawdata, &data)
					if err != nil {
						return err
					}
					parsedTime, err := time.Parse(time.RFC3339, string(timestamp))
					if err != nil {
						return err
					}
					entries = append(entries, Entry{
						Collection: strName,
						Timestamp:  parsedTime,
						Data:       data,
					})
				}
			}

			return nil
		})
	})
	if err != nil {
		return []Entry{}, err
	}

	return entries, nil
}

func (app *App) GetLatestFromDB() ([]Entry, error) {
	var entries []Entry
	err := app.DB.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, bucket *bolt.Bucket) error {
			cursor := bucket.Cursor()

			timestamp, rawdata := cursor.Last()

			var data map[string]interface{}
			err := json.Unmarshal(rawdata, &data)
			if err != nil {
				return err
			}
			parsedTime, err := time.Parse(time.RFC3339, string(timestamp))
			if err != nil {
				return err
			}
			entries = append(entries, Entry{
				Collection: string(name),
				Timestamp:  parsedTime,
				Data:       data,
			})

			return nil
		})
	})
	if err != nil {
		return []Entry{}, err
	}

	return entries, nil
}

func (app *App) CreateEntryInDB(entry Entry) error {
	insertableEntry, err := ValidateAndConvertEntry(entry)
	if err != nil {
		return err
	}
	err = app.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(insertableEntry.Collection))
		if err != nil {
			return err
		}
		return b.Put(insertableEntry.Timestamp, insertableEntry.Data)
	})
	if err != nil {
		return err
	}

	return nil
}
