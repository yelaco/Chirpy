package database

import (
	"errors"
)

type User struct {
	Id       int
	Password string
	Email    string
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateUser(password string, email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("CreateUser: " + err.Error())
	}

	hashedPassword, err := getHashedPassword(password)
	if err != nil {
		return User{}, errors.New("CreateUser: " + err.Error())
	}
	user := User{
		Id:       len(dbs.Users) + 1,
		Password: hashedPassword,
		Email:    email,
	}

	dbs.Users[user.Id] = user
	err = db.writeDB(dbs)
	if err != nil {
		return User{}, errors.New("CreateUser: " + err.Error())
	}

	return user, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetUserByEmail(email string) (User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("GetUserByEmail: " + err.Error())
	}

	for _, user := range dbs.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("Could not find user")
}

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
