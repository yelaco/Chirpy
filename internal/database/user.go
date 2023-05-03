package database

import (
	"errors"
	"strconv"
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

	hashedPassword, err := GetHashedPassword(password)
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

func (db *DB) GetUserById(userId string) (User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("GetUserById: " + err.Error())
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return User{}, errors.New("GetUserById: " + err.Error())
	}
	if user, ok := dbs.Users[id]; ok {
		return user, nil
	}
	return User{}, errors.New("Not found")
}

func (db *DB) UpdateUser(userId string, password string, email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return User{}, errors.New("UpdateUser: " + err.Error())
	}

	hashedPassword, err := GetHashedPassword(password)
	if err != nil {
		return User{}, errors.New("UpdateUser: " + err.Error())
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return User{}, errors.New("UpdateUser: " + err.Error())
	}

	if _, ok := dbs.Users[id]; ok {
		dbs.Users[id] = User{
			Id:       id,
			Password: hashedPassword,
			Email:    email,
		}
	} else {
		return User{}, errors.New("UpdateUser: User not found")
	}

	err = db.writeDB(dbs)
	if err != nil {
		return User{}, errors.New("UpdateUser: " + err.Error())
	}

	return dbs.Users[id], nil
}
