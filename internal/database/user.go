package database

import (
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Password string
	Email    string
}

func GetHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("Failed to hash password")
	}
	return string(hashedPassword), nil
}

func (db *DB) CreateUser(password string, email string) (*User, error) {
	if usr, _ := db.GetUserByEmail(email); usr != nil {
		return nil, errors.New("CreateUser: User already exists")
	}

	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("CreateUser: " + err.Error())
	}

	hashedPassword, err := GetHashedPassword(password)
	if err != nil {
		return nil, errors.New("CreateUser: " + err.Error())
	}
	user := &User{
		Id:       len(dbs.Users) + 1,
		Password: hashedPassword,
		Email:    email,
	}

	dbs.Users[user.Id] = *user
	err = db.writeDB(dbs)
	if err != nil {
		return nil, errors.New("CreateUser: " + err.Error())
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (*User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("GetUserByEmail: " + err.Error())
	}

	for _, user := range dbs.Users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("Could not find user")
}

func (db *DB) GetUserById(userId string) (*User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("GetUserById: " + err.Error())
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, errors.New("GetUserById: " + err.Error())
	}
	if user, ok := dbs.Users[id]; ok {
		return &user, nil
	}
	return nil, errors.New("Not found")
}

func (db *DB) UpdateUser(userId string, password string, email string) (*User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("UpdateUser: " + err.Error())
	}

	hashedPassword, err := GetHashedPassword(password)
	if err != nil {
		return nil, errors.New("UpdateUser: " + err.Error())
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, errors.New("UpdateUser: " + err.Error())
	}

	if _, ok := dbs.Users[id]; ok {
		dbs.Users[id] = User{
			Id:       id,
			Password: hashedPassword,
			Email:    email,
		}
	} else {
		return nil, errors.New("UpdateUser: User not found")
	}

	err = db.writeDB(dbs)
	if err != nil {
		return nil, errors.New("UpdateUser: " + err.Error())
	}

	user := dbs.Users[id]
	return &user, nil
}
