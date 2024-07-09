package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotExist = errors.New("user not exist")

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	IsChirpyRed    bool   `json:"is_chirpy_red"`
}

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
		IsChirpyRed:    false,
	}
	dbStructure.Users[id] = user

	if err := db.writeDB(dbStructure); err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeChirpyRed(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	user, err := db.GetUserByID(id)
	if err != nil {
		return err
	}

	user.IsChirpyRed = true
	dbStructure.Users[id] = user

	if err := db.writeDB(dbStructure); err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	if _, err := db.GetUserByID(id); err != nil {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return User{}, err
	}

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

	return User{}, ErrUserNotExist
}

func (db *DB) GetUserByID(userID int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	users := dbStructure.Users
	for _, user := range users {
		if user.ID == userID {
			return user, nil
		}
	}

	return User{}, ErrUserNotExist
}
