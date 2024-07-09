package main

import (
	"encoding/json"
	"net/http"

	"github.com/yelaco/Chirpy/internal/auth"
	"github.com/yelaco/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string         `json:"event"`
		Data  map[string]int `json:"data"`
	}

	if apiKey, err := auth.GetApiKey(r.Header); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
		return
	} else if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid api key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Event ignored")
		return
	}

	userID, exists := params.Data["user_id"]
	if !exists {
		respondWithError(w, http.StatusBadRequest, "Couldn't find user ID")
	}

	err := cfg.DB.UpgradeChirpyRed(userID)
	if err == database.ErrUserNotExist {
		respondWithError(w, http.StatusNotFound, "Coudn't find user")
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
