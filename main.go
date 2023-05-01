package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	apiCfg := &apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()

	// Register handlers and its according request methods to the api router
	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", handleReadiness)
	apiRouter.Get("/metrics", apiCfg.handlerMetrics)

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
