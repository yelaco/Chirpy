package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
	Password    string `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	newUser, err := cfg.DB.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:          newUser.ID,
			Email:       newUser.Email,
			IsChirpyRed: newUser.IsChirpyRed,
		},
	})
}
