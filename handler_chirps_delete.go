package main

import (
	"net/http"
	"strconv"

	"github.com/yelaco/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find access token")
		return
	}

	userIDString, err := auth.ValidateAccessToken(accessToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate accesss token")
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid access token")
		return
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirpid"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirpFromID(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Coulnd't find chirp")
		return
	}

	if chirp.AuthorID != userID {
		respondWithError(w, http.StatusForbidden, "Forbidden operation")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "Chirp deleted")
}
