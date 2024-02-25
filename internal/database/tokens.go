package database

import (
	"errors"
	"time"
)

func (db *DB) RevokeToken(t string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return errors.New("Failed to connect to DB")
	}

	dbStructure.Tokens[t] = time.Now() 
	return db.writeDB(dbStructure)
}

func (db *DB) IsTokenRevoked(t string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, errors.New("Failed to connect to DB")
	}

	_, ok := dbStructure.Tokens[t] 
	return ok, nil
} 

