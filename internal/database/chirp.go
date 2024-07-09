package database

import (
	"log"
	"os"
)

type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(userID int, body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		AuthorID: userID,
		Body:     body,
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

func (db *DB) GetChirpFromID(id int) (Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		log.Printf("couldn't load database - %v", err)
		return Chirp{}, err
	}

	chirp, exist := dbs.Chirps[id]
	if !exist {
		return Chirp{}, os.ErrNotExist
	}

	return chirp, nil
}

func (db *DB) GetChirpsFromAuthorID(authorID int, order string) ([]Chirp, error) {
	dbs, err := db.loadDB()
	if err != nil {
		log.Printf("couldn't load database - %v", err)
		return []Chirp{}, err
	}

	chirps := []Chirp{}

	for _, chirp := range dbs.Chirps {
		if chirp.AuthorID == authorID {
			chirps = append(chirps, chirp)
		}
	}

	return chirps, nil
}
