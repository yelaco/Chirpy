package database

import (
	"encoding/json"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		err = os.WriteFile(db.path, []byte(""), 0666)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs := DBStructure{
		Chirps: map[int]Chirp{},
		Users:  map[int]User{},
	}

	file, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		return DBStructure{}, os.ErrNotExist
	}

	if len(file) > 0 {
		if err := json.Unmarshal(file, &dbs); err != nil {
			return DBStructure{}, err
		}
	}

	return dbs, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.path, dat, 0666); err != nil {
		return err
	}

	return nil
}

// delete the existing database file and create a new one on disk
func (db *DB) ResetDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	if err := os.Remove(db.path); err != nil {
		return err
	}

	if err := db.ensureDB(); err != nil {
		return err
	}

	return nil
}
