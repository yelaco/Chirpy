package database

import (
	"errors"
	"strconv"
)

type Chirp struct {
	Id        int
	Author_Id int
	Body      string
}

func getMaxId(chirps map[int]Chirp) int {
	var maxId int
	for maxId = range chirps {
	}
	return maxId
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, errors.New("CreateChirp: " + err.Error())
	}

	chirp := Chirp{
		Id:        getMaxId(dbs.Chirps) + 1,
		Author_Id: authorId,
		Body:      body,
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

func (db *DB) DeleteChirpById(chirpId string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return errors.New("DeleteChirpById: " + err.Error())
	}

	id, err := strconv.Atoi(chirpId)
	if err != nil {
		return errors.New("DeleteChirpById: " + err.Error())
	}
	for k, v := range dbs.Chirps {
		if v.Id == id {
			delete(dbs.Chirps, k)
			err = db.writeDB(dbs)
			if err != nil {
				return errors.New("DeleteChirpById: " + err.Error())
			}
			return nil
		}
	}
	return errors.New("DeleteChirpById: Id not found")
}
