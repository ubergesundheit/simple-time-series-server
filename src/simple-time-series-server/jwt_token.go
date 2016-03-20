package main

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	MissingCollectionError = errors.New("Key `collection` must be present and non empty")
	MissingTimestampError  = errors.New("Key `timestamp` must be present and non empty")
	MissingDataError       = errors.New("Key `data` must be present and non empty")
)

func ParseJwt(rawJwt string) (Entry, error) {
	createdAt, err := time.Now().UTC().MarshalText()
	if err != nil {
		return Entry{}, err
	}
	token, err := jwt.Parse(rawJwt, func(token *jwt.Token) (interface{}, error) {
		if collection, ok := token.Claims["collection"]; ok {
			if collectionStr, ok := collection.(string); ok {
				return CollectionSecrets[collectionStr], nil
			} else {
				return 0, MissingCollectionError
			}
		} else {
			return 0, MissingCollectionError
		}

	})
	if err != nil {
		return Entry{}, err
	} else if token.Valid {
		var newJwtEntry = Entry{}
		if collection, ok := token.Claims["collection"].(string); ok {
			if len(collection) == 0 {
				return Entry{}, MissingCollectionError
			}
			newJwtEntry.Collection = collection
		} else {
			return Entry{}, MissingCollectionError
		}

		if ts, ok := token.Claims["timestamp"].(string); ok {
			if len(ts) == 0 {
				return Entry{}, MissingTimestampError
			}
			// try to parse the timestamp
			parsedTime, err := time.Parse(time.RFC3339, ts)
			if err != nil {
				return Entry{}, err
			}

			newJwtEntry.Timestamp = parsedTime
		} else {
			return Entry{}, MissingTimestampError
		}

		if data, ok := token.Claims["data"].(map[string]interface{}); ok {
			if len(data) == 0 {
				return Entry{}, MissingDataError
			}

			newJwtEntry.Data = data
			newJwtEntry.Data["createdAt"] = string(createdAt)
		} else {
			return Entry{}, MissingDataError
		}

		return newJwtEntry, nil
	}

	return Entry{}, nil
}

//func validateJwtSignature(header, payload, signature) error {
//	mac := hmac.New(sha256.New, SECRET)

//	mac.Write(strings.Join(header, payload, "."))

//}

//func validateJwtHeader(headerStr string) error {
//	// for now, just compare to {"alg":"HS256","typ":"JWT"}
//	header, err := base64.StdEncoding.DecodeString(headerStr)
//	if err != nil {
//		return err
//	}
//	if bytes.Compare(header, ValidHeader) != 0 {
//		return WrongJwtHeaderError
//	}
//	return nil
//}

//func unpackJwt(rawJwt []byte) (string, string, string) {
//	parts := strings.Split(string(rawJwt), ".")
//	return parts[0], parts[1], parts[2]
//}
