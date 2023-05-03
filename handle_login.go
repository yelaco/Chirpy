package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cf *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password           string
		Email              string
		Expires_In_Seconds int
		Test_You_Here      string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handleLogin: " + err.Error())
		return
	}

	expireAt := params.Expires_In_Seconds

	if loginResp, valid := cf.validatePassword(params.Password, params.Email); valid {
		defaultExpiration := 60 * 60 * 24
		if expireAt == 0 {
			expireAt = defaultExpiration
		} else if expireAt > defaultExpiration {
			expireAt = defaultExpiration
		}
		loginResp.Token, err = cf.getSignedToken(expireAt, loginResp.Id)
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
