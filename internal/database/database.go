package database 

import (
	"encoding/json"
	"errors"
	"time"
	"os"
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User `json:"users"`
	Tokens map[string]time.Time `json:"tokens"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu: &sync.RWMutex{},
	} 
	err := db.ensureDB()
	return db, err
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
		Tokens: map[string]time.Time{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	dbStructure := DBStructure{}
	j, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	} 	

	err = json.Unmarshal(j, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	j, err := json.Marshal(dbStructure) 
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, []byte(j), 0666)
	if err != nil {
		return err
	}
	return nil
}

