package main

import "errors"

func (cf *apiConfig) authCreateChirp(accessToken string) (int, error) {
	if valid, err := cf.validateAccessToken(accessToken); !valid {
		return -1, errors.New("authCreateChirp: " + err.Error())
	}
	userId, err := cf.userIdFromToken(accessToken)
	if err != nil {
		return -1, errors.New("authCreateChirp: " + err.Error())
	}
	return userId, nil
}

func (cf *apiConfig) authDeleteChirp(accessToken string, authorId int) (bool, error) {
	if valid, err := cf.validateAccessToken(accessToken); !valid {
		return false, errors.New("authDeleteChirp: " + err.Error())
	}
	userId, err := cf.userIdFromToken(accessToken)
	if err != nil {
		return false, errors.New("authDeleteChirp: " + err.Error())
	}
	if authorId != userId {
		return false, errors.New("authDeleteChirp: Unauthorized user")
	}
	return true, nil
}
