package main

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (cf *apiConfig) validatePassword(password string, email string) (UserResponse, bool) {
	user, err := cf.db.GetUserByEmail(email)
	if err != nil {
		return UserResponse{}, false
	}
	e := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if e != nil {
		return UserResponse{}, false
	}
	return UserResponse{user.Id, user.Email}, true
}

func (cf *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string
		Email    string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handleLogin: " + err.Error())
		return
	}
	if userResp, valid := cf.validatePassword(params.Password, params.Email); valid {
		respondWithJSON(w, http.StatusOK, userResp)
		return
	}
	respondWithError(w, http.StatusUnauthorized, "Password does not match")
}
