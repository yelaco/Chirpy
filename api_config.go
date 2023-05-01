package main

import "github.com/minhquang053/Chirpy/internal/database"

type apiConfig struct {
	db             *database.DB
	fileserverHits int
}
