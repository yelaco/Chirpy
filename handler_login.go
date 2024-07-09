package main

import (
	"encoding/json"
	"net/http"

	"github.com/yelaco/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find user with specified email")
		return
	}

	if err := auth.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password not matched")
		return
	}

	defaultExpiration := 60 * 60
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	accessToken, err := auth.CreateAccessToken(user.ID, params.ExpiresInSeconds, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coulnd't create access token")
		return
	}

	refreshToken, err := auth.CreateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coulnd't create refresh token")
		return
	}

	if _, err := cfg.DB.CreateTokenInfo(user.ID, refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coulnd't create refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:           user.ID,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
