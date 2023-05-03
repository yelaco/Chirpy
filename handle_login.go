package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cf *apiConfig) validatePassword(password string, email string) (LoginResponse, bool) {
	login, err := cf.db.GetUserByEmail(email)
	if err != nil {
		return LoginResponse{}, false
	}
	e := bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(password))
	if e != nil {
		return LoginResponse{}, false
	}
	return LoginResponse{
		Id:    login.Id,
		Email: login.Email,
	}, true
}

func (cf *apiConfig) getSignedToken(expireInSecond int, userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(expireInSecond) * time.Second)),
		Subject:   fmt.Sprintf("%v", userId),
	})
	signedToken, err := token.SignedString([]byte(cf.jwtSecret))
	if err != nil {
		return "", errors.New("getSignedToken: " + err.Error())
	}
	return signedToken, nil
}

func (cf *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password          string
		Email             string
		Expire_In_Seconds int
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handleLogin: " + err.Error())
		return
	}
	if loginResp, valid := cf.validatePassword(params.Password, params.Email); valid {
		loginResp.Token, err = cf.getSignedToken(params.Expire_In_Seconds, loginResp.Id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't provide token")
			log.Println("handleLogin: " + err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, loginResp)
		return
	}
	respondWithError(w, http.StatusUnauthorized, "Password does not match")
}
