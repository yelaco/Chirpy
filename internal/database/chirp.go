package database

import (
	"log"
	"os"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		log.Printf("GetChirps(): couldn't load database - %v", err)
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbs.Chirps))
	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirpFromID(id int) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		log.Printf("GetChirps(): couldn't load database - %v", err)
		return Chirp{}, err
	}

	chirp, exist := dbs.Chirps[id]
	if !exist {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}
