package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	filenameDefault = "./simple-time-series-db.db"
	filenameUsage   = "filename of the boltdb database"
	addressDefault  = "127.0.0.1:8080"
	addressUsage    = "adress to bind to"
	dropDefault     = false
	dropUsage       = "only for import command: Drop the database before importing"
	secretsDefault  = "./simple-time-series-secrets"
	secretsUsage    = "specify the filename of JWT secrets, file must contain lines in format <collection>=<secret>"
)

var (
	Filename       string
	Address        string
	DropDB         bool
	SecretFilename string

	CollectionSecrets map[string][]byte = make(map[string][]byte)

	MalformedSecretsFileError = func(linenum int, msg string) error {
		return errors.New("Malformed secrets file at line " + strconv.Itoa(linenum) + " " + msg)
	}
)

func handleConfig() {
	flag.StringVar(&Filename, "filename", filenameDefault, filenameUsage)
	flag.StringVar(&Filename, "f", filenameDefault, filenameUsage)

	flag.StringVar(&Address, "address", addressDefault, addressUsage)
	flag.StringVar(&Address, "A", addressDefault, addressUsage)

	flag.BoolVar(&DropDB, "drop", dropDefault, dropUsage)

	flag.StringVar(&SecretFilename, "secrets", secretsDefault, secretsUsage)
	flag.StringVar(&SecretFilename, "s", secretsDefault, secretsUsage)

	flag.Parse()

	err := parseSecrets(CollectionSecrets)
	if err != nil {
		fmt.Println("No secrets available! Error:", err)
	}
}

func parseSecrets(secrets map[string][]byte) error {
	file, err := os.Open(SecretFilename)
	defer file.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// read lines and split them by the equal sign (=)
	linenum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 { // ignore empty lines
			keyValue := strings.SplitN(line, "=", 2)
			if len(keyValue) != 2 || len(keyValue[0]) == 0 || len(keyValue[1]) == 0 {
				return MalformedSecretsFileError(linenum, "")
			} else if len(secrets[keyValue[0]]) != 0 {
				return MalformedSecretsFileError(linenum, "duplicate key")
			} else {
				secrets[keyValue[0]] = []byte(keyValue[1])
			}
		}
		linenum += 1
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
