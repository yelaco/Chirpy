package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/minhquang053/Chirpy/internal/database"
)

func main() {
	// Look for a file name .env in current directory to load environment variables
	godotenv.Load()

	apiCfg := &apiConfig{
		db:             database.NewDB("internal/database/database.json"),
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	r := chi.NewRouter()

	// Register handlers and its according request methods to the api router
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)
	apiRouter.Post("/chirps", apiCfg.handlePostChirp)
	apiRouter.Get("/chirps", apiCfg.handleGetAllChirps)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handleGetChirpById)
	apiRouter.Post("/users", apiCfg.handlePostUser)
	apiRouter.Put("/users", apiCfg.handlePutUser)
	apiRouter.Post("/login", apiCfg.handleLogin)
	apiRouter.Post("/refresh", apiCfg.handleRefresh)

	// Regiser handlers and its according request methods to the admin router
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.adminHandlerMetrics)

	// Mount the subrouters -> /api/healthz will be handled by apiRouter with /healthz pattern
	r.Mount("/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	r.Mount("/api", apiRouter)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCors(r)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	log.Fatal(server.ListenAndServe())
}
