package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	dat, err := json.Marshal(struct{ Error string }{msg})

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func replaceBadWords(s string) string {
	words := strings.Split(s, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if loweredWord == "kerfuffle" || loweredWord == "sharbert" || loweredWord == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cf *apiConfig) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	modifiedBody := replaceBadWords(params.Body)

	newChirp, err := cf.db.CreateChirp(modifiedBody)
	if err != nil {
		log.Println("handlePostChirp: " + err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, newChirp)
}

func (cf *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	type response []struct {
		Id   int
		Body string
	}
	chirps, err := cf.db.GetChirps()
	resp := response{}
	for _, ch := range chirps {
		resp = append(resp, struct {
			Id   int
			Body string
		}{ch.Id, ch.Body})
	}
	if err != nil {
		log.Println("handleGetChirp: " + err.Error())
		return
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Id < resp[j].Id
	})

	respondWithJSON(w, http.StatusOK, resp)
}
