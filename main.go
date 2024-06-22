package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/yelaco/Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const (
		port         string = "8080"
		filepathRoot string = "."
	)
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}
	if *dbg {
		db.ResetDB()
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpid}", apiCfg.handlerChirpsGetFromID)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
