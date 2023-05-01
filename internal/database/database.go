package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) *DB {
	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp := Chirp{
		Id:   len(dbs.Chirps) + 1,
		Body: body,
	}

	dbs.Chirps[chirp.Id] = chirp
	err = db.writeDB(dbs)
	if err != nil {
		return Chirp{}, err
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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		f, err := os.Create(db.path)
		if err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)

	if err == os.ErrNotExist {
		e := db.ensureDB()
		if e != nil {
			return DBStructure{map[int]Chirp{}}, e
		}
		return db.loadDB()
	}

	if len(data) == 0 {
		return DBStructure{map[int]Chirp{}}, nil
	}

	dbs := DBStructure{map[int]Chirp{}}
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return DBStructure{map[int]Chirp{}}, errors.New("loadDB error: " + err.Error())
	}

	return dbs, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return errors.New("writeDB error : " + err.Error())
	}

	return os.WriteFile(db.path, data, 0644)
}
