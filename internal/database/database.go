package database

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp
	Users         map[int]User
	RefreshTokens map[string]RefreshToken
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) *DB {
	// will delete the database file at the start of the program if in debug mode
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		deleteDB(path)
	}

	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
}

func emptyDBS() DBStructure {
	return DBStructure{
		map[int]Chirp{},
		map[int]User{},
		map[string]RefreshToken{},
	}
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		f, err := os.Create(db.path)
		if err != nil {
			return errors.New("ensureDB: " + err.Error())
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
			return emptyDBS(), errors.New("loadDB: " + err.Error())
		}
		return db.loadDB()
	}

	if len(data) == 0 {
		return emptyDBS(), nil
	}

	dbs := emptyDBS()
	err = json.Unmarshal(data, &dbs)
	if err != nil {
		return emptyDBS(), errors.New("loadDB: " + err.Error())
	}

	return dbs, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return errors.New("writeDB: " + err.Error())
	}

	return os.WriteFile(db.path, data, 0644)
}

// Delete database json file (for debugging purposes only)
func deleteDB(path string) error {
	return os.Remove(path)
}
