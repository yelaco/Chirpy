package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	mux.Handle("/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/healthz", ReadinessEndpoint)

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	log.Fatal(server.ListenAndServe())
}
