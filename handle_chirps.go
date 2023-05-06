package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/minhquang053/Chirpy/internal/database"
)

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
	accessToken, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		log.Println("handlePostChirp: " + err.Error())
		return
	}

	authorId, err := cf.authCreateChirp(accessToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		log.Println("handlePostChirp: " + err.Error())
		return
	}

	type parameters struct {
		Body string
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode JSON")
		log.Println("handlePostChirp: " + err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	modifiedBody := replaceBadWords(params.Body)

	newChirp, err := cf.db.CreateChirp(modifiedBody, authorId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp")
		log.Println("handlePostChirp: " + err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, newChirp)
}

func (cf *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	type response []struct {
		Id        int
		Author_Id int
		Body      string
	}
	var chirps []database.Chirp

	authorId := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	if len(authorId) == 0 {
		ch, err := cf.db.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps")
			log.Println("handleGetChirps: " + err.Error())
			return
		}
		chirps = ch
	} else {
		ch, err := cf.db.GetChirpsByAuthorId(authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirp")
			log.Println("handleGetChirps: " + err.Error())
			return
		}
		chirps = ch
	}

	resp := response{}
	for _, ch := range chirps {
		resp = append(resp, struct {
			Id        int
			Author_Id int
			Body      string
		}{ch.Id, ch.Author_Id, ch.Body})
	}

	sort.Slice(resp, func(sortOrder string) func(int, int) bool {
		// Return function to sort the slice in order specified
		switch sortOrder {
		case "asc":
			return func(i, j int) bool {
				return resp[i].Id < resp[j].Id
			}
		case "desc":
			return func(i, j int) bool {
				return resp[i].Id > resp[j].Id
			}
		default:
			return func(i, j int) bool {
				return resp[i].Id < resp[j].Id
			}
		}
	}(sortOrder))

	respondWithJSON(w, http.StatusOK, resp)
}

func (cf *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
	chirpId := chi.URLParam(r, "chirpID")
	chirp, err := cf.db.GetChirpById(chirpId)
	if err != nil {
		if err.Error() == "Not found" {
			respondWithError(w, http.StatusNotFound, "Id not found")
			return
		}
		log.Println("handleGetChirpById: " + err.Error())
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirp")
	}
	respondWithJSON(w, http.StatusOK, chirp)
}

func (cf *apiConfig) handleDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	accessToken, err := GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		log.Println("handlePostChirp: " + err.Error())
		return
	}

	chirpId := chi.URLParam(r, "chirpID")
	chirp, err := cf.db.GetChirpById(chirpId)
	if err != nil {
		if err.Error() == "Not found" {
			respondWithError(w, http.StatusNotFound, "Id not found")
			return
		}
		log.Println("handleDeleteChirpById: " + err.Error())
		respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirp")
	}
	if valid, err := cf.authDeleteChirp(accessToken, chirp.Author_Id); !valid {
		log.Println("handleDeleteChirpById: " + err.Error())
		respondWithError(w, http.StatusForbidden, "Unauthorized user")
		return
	}
	err = cf.db.DeleteChirpById(chirpId)
	if err != nil {
		log.Println("handleDeleteChirpById: " + err.Error())
	}
	respondEmpty(w, http.StatusOK)
}
