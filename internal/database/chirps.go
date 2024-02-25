package database

import "os"

type Chirp struct {
	ID   int `json:"id"`	
	Body string `json:"body"`
	AuthorID int `json:"author_id"`
}


func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	dbStructure, err := db.loadDB()	
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID: id,
		Body: body,
		AuthorID: authorId,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}

	dbStructure, err := db.loadDB()	
	if err != nil {
		return chirps, err
	}

	for _, c := range dbStructure.Chirps {
		chirps = append(chirps, c)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}

