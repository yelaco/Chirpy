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
	corsMux := middlewareCors(r)
	r.Mount("/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))

	r.Get("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	r.Get("/metrics", apiCfg.handlerMetrics)
	// mux := http.NewServeMux()
	// corsMux := middlewareCors(mux)
	// apiCfg := &apiConfig{0}

	// mux.Handle("/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	// mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("OK"))
	// })
	// mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	log.Fatal(server.ListenAndServe())
}
