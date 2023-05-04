package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserResponse struct {
	Id    int
	Email string
}

type LoginResponse struct {
	Id            int
	Email         string
	Token         string
	Refresh_Token string
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("respondWithJSON: " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}
