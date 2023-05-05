package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

	if loginResp, valid := cf.validatePassword(params.Password, params.Email); valid {
		loginResp.Token, err = cf.createAccessToken(loginResp.Id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't provide token")
			log.Println("handleLogin: " + err.Error())
			return
		}
		loginResp.Refresh_Token, err = cf.createRefreshToken(loginResp.Id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't provide refresh token")
			log.Println("handleLogin: " + err.Error())
			return
		}
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't add refresh token")
			log.Println("handleLogin: " + err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, loginResp)
		return
	}
	respondWithError(w, http.StatusUnauthorized, "Password does not match")
}
