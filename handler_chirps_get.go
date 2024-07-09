package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/yelaco/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	var dbChirps []database.Chirp

	order := r.URL.Query().Get("sort")

	authorID, err := strconv.Atoi(r.URL.Query().Get("author_id"))
	if err == nil {
		dbChirps, err = cfg.DB.GetChirpsFromAuthorID(authorID, order)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirps")
			return
		}
	} else {
		dbChirps, err = cfg.DB.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirps")
			return
		}
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:       dbChirp.ID,
			AuthorID: dbChirp.AuthorID,
			Body:     dbChirp.Body,
		})
	}

	if order == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsGetFromID(w http.ResponseWriter, r *http.Request) {
	chirpId, err := strconv.Atoi(r.PathValue("chirpid"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirpFromID(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}
