package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateUser(email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if _, err := db.GetUserByEmail(email); err == nil {
		return User{}, errors.New("user existed")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:             id,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}
	dbStructure.Users[id] = user

	if err := db.writeDB(dbStructure); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	users := dbStructure.Users
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user not exist")
}
