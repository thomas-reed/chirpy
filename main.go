package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"
	const filepathRoot = "."
	
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(fsHandler))
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /metrics", apiConfig.metricsHandler)
	mux.HandleFunc("POST /reset", apiConfig.resetHandler)

	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Server listening on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}