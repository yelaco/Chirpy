package database

import (
	"errors"
	"time"
)

type TokenInfo struct {
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) CreateTokenInfo(userID int, tokenString string) (TokenInfo, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return TokenInfo{}, err
	}

	token := TokenInfo{
		UserID:    userID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	}
	dbStructure.TokenInfos[tokenString] = token

	if err := db.writeDB(dbStructure); err != nil {
		return TokenInfo{}, err
	}

	return token, nil
}

func (db *DB) GetTokenInfo(tokenString string) (TokenInfo, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return TokenInfo{}, err
	}

	if token, ok := dbStructure.TokenInfos[tokenString]; ok {
		return token, nil
	}

	return TokenInfo{}, errors.New("refresh token not exist")
}

func (db *DB) DeleteTokenInfo(tokenString string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.TokenInfos, tokenString)

	return db.writeDB(dbStructure)
}
