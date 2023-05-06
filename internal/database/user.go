package database

import (
	"errors"
	"fmt"
	"strconv"
)

type User struct {
	Id            int
	Password      string
	Email         string
	Is_Chirpy_Red bool
}

func (db *DB) CreateUser(password *string, email string) (*User, error) {
	if password == nil {
		return nil, errors.New("Couldn't hash password")
	}

	if usr, _ := db.GetUserByEmail(email); usr != nil {
		return nil, errors.New("CreateUser: User already exists")
	}

	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("CreateUser: " + err.Error())
	}

	user := &User{
		Id:            len(dbs.Users) + 1,
		Password:      *password,
		Email:         email,
		Is_Chirpy_Red: false,
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

func (db *DB) UpdateUser(userId string, password *string, email string, isRed bool) (*User, error) {
	if password == nil {
		return nil, errors.New("Couldn't hash password")
	}

	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("UpdateUser: " + err.Error())
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return nil, errors.New("UpdateUser: " + err.Error())
	}

	if _, ok := dbs.Users[id]; ok {
		dbs.Users[id] = User{
			Id:            id,
			Password:      *password,
			Email:         email,
			Is_Chirpy_Red: isRed,
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

func (db *DB) UpgradeUser(id int) error {
	userId := fmt.Sprintf("%d", id)
	user, err := db.GetUserById(userId)
	if user == nil {
		return errors.New("UpgradeUser: " + err.Error())
	}
	_, e := db.UpdateUser(userId, &user.Password, user.Email, true)
	if e != nil {
		return errors.New("UpgradeUser: " + e.Error())
	}
	return nil
}
