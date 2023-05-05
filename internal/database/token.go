package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	RevokeTime time.Time
}

func (db *DB) RevokeRefreshToken(signedToken string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbs, err := db.loadDB()
	if err != nil {
		return errors.New("RevokeRefreshToken: " + err.Error())
	}

	dbs.RefreshTokens[signedToken] = RefreshToken{
		RevokeTime: time.Now(),
	}
	err = db.writeDB(dbs)
	if err != nil {
		return errors.New("RevokeRefreshToken: " + err.Error())
	}

	return nil
}

func (db *DB) GetRefreshToken(refreshToken string) (*RefreshToken, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbs, err := db.loadDB()
	if err != nil {
		return nil, errors.New("GetRefreshToken " + err.Error())
	}

	if token, ok := dbs.RefreshTokens[refreshToken]; ok {
		return &token, nil
	}
	return nil, errors.New("Not found")
}
