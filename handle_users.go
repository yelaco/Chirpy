package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cf *apiConfig) handlePostUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handlePostUser: " + err.Error())
		return
	}

	newUser, err := cf.db.CreateUser(params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create User")
		log.Println("handlePostUser: " + err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, newUser)
}
