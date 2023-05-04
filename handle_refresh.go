package main

import (
	"log"
	"net/http"
)

func (cf *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad refresh token")
		log.Println("handleRefresh: " + err.Error())
		return
	}

	if valid, err := cf.isRefreshToken(refreshToken); !valid {
		respondWithError(w, http.StatusUnauthorized, "Bad refresh token")
		log.Println("handleRefresh: " + err.Error())
		return
	}

	userId, err := cf.userIdFromToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bad refresh token")
		log.Println("handleRefresh: " + err.Error())
		return
	}

	newRefreshToken, err := cf.getRefreshToken(userId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not get refresh token")
	}

	respondWithJSON(w, http.StatusOK, struct{ Token string }{newRefreshToken})
}
