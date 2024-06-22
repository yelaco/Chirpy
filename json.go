package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type ChirpResponse struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
	w.WriteHeader(code)
	w.Write(dat)
}
