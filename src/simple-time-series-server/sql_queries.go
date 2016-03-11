package main

const (
	INSERT_INTO             string = "INSERT INTO `entries` (collection, timestamp, data) VALUES (?, ?, ?);"
	SELECT_LATEST           string = "SELECT `collection`, MAX(`timestamp`) as timestamp, `data` FROM `entries` GROUP BY collection;"
	CREATE_TABLE            string = "CREATE TABLE IF NOT EXISTS entries (collection TEXT, timestamp INTEGER, data BLOB);"
	CREATE_INDEX_COLLECTION string = "CREATE INDEX IF NOT EXISTS index_collection ON entries (collection);"
	PRAGMAS                 string = "PRAGMA locking_mode = EXCLUSIVE;PRAGMA synchronous = OFF;PRAGMA journal_mode = OFF;"
)
