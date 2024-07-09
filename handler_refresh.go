package main

import (
	"net/http"

	"github.com/yelaco/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
		return
	}

	tokenInfo, err := cfg.DB.GetTokenInfo(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	userID, err := auth.ValidateRefreshToken(tokenInfo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse user ID")
		return
	}

	defaultExpiration := 60 * 60
	accessToken, err := auth.CreateAccessToken(userID, defaultExpiration, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coulnd't create access token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
		return
	}

	if err := cfg.DB.DeleteTokenInfo(tokenString); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
