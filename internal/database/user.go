package database

import (
	"errors"
)

type User struct {
	Id    int
	Email string
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateUser(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("CreateUser: " + err.Error())
	}

	user := User{
		Id:    len(dbs.Users) + 1,
		Email: email,
	}

	dbs.Users[user.Id] = user
	err = db.writeDB(dbs)
	if err != nil {
		return User{}, errors.New("CreateUser: " + err.Error())
	}

	return user, nil
}

// // GetChirps returns all chirps in the database
// func (db *DB) GetUsers() ([]Chirp, error) {
// 	db.mux.RLock()
// 	defer db.mux.RUnlock()

// 	dbs, err := db.loadDB()
// 	if err != nil {
// 		return nil, errors.New("GetChirps: " + err.Error())
// 	}
// 	chirps := []Chirp{}
// 	for _, chirp := range dbs.Chirps {
// 		chirps = append(chirps, chirp)
// 	}

// 	return chirps, nil
// }

// func (db *DB) GetChirpById(chirpId string) (Chirp, error) {
// 	db.mux.RLock()
// 	defer db.mux.RUnlock()

// 	dbs, err := db.loadDB()
// 	if err != nil {
// 		return Chirp{}, errors.New("GetChirpById: " + err.Error())
// 	}

// 	id, err := strconv.Atoi(chirpId)
// 	if err != nil {
// 		return Chirp{}, errors.New("GetChirpById: " + err.Error())
// 	}
// 	if chirp, ok := dbs.Chirps[id]; ok {
// 		return chirp, nil
// 	}
// 	return Chirp{}, errors.New("Not found")
// }
