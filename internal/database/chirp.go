package database

import (
	"errors"
	"strconv"
)

type Chirp struct {
	Id   int
	Body string
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, errors.New("CreateChirp: " + err.Error())
	}

	chirp := Chirp{
		Id:   len(dbs.Chirps) + 1,
		Body: body,
	}

	dbs.Chirps[chirp.Id] = chirp
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, errors.New("CreateChirp: " + err.Error())
	}

	return chirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("GetChirps: " + err.Error())
	}
	chirps := []Chirp{}
	for _, chirp := range dbs.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirpById(chirpId string) (Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, errors.New("GetChirpById: " + err.Error())
	}

	id, err := strconv.Atoi(chirpId)
	if err != nil {
		return Chirp{}, errors.New("GetChirpById: " + err.Error())
	}
	if chirp, ok := dbs.Chirps[id]; ok {
		return chirp, nil
	}
	return Chirp{}, errors.New("Not found")
}
