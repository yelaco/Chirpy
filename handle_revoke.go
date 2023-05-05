package main

import "net/http"

func (cf *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not retrieve token")
		return
	}
	err = cf.db.RevokeRefreshToken(token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not revoke refresh token")
		return
	}

	respondWithSuccess(w, http.StatusOK, "Revoke successfully")
}
