package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cf *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, e := GetApiKey(r.Header)
	if e != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not retrieve api key from headers")
		log.Println("handleWebhook: " + e.Error())
		return
	}
	if apiKey != cf.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized api key")
		return
	}

	type parameters struct {
		Event string         `json:"event"`
		Data  map[string]int `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handleLogin: " + err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusOK, "Invalid event")
		log.Println("handleWebhook: Received invalid event")
		return
	}

	err = cf.db.UpgradeUser(params.Data["user_id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		log.Println("handleWebhook: " + err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
