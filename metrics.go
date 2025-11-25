package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		`
			<html>
				<body>
					<h1>Welcome, Chirpy Admin</h1>
					<p>Chirpy has been visited %d times!</p>
				</body>
			</html>
		`,
		cfg.fileserverHits.Load(),
	)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	returnStr := "Hits reset to 0"
	if cfg.platform == "dev" {
		ctx := context.Background()
		err := cfg.db.DeleteAllUsers(ctx)
		if err != nil {
			log.Printf("Error Deleting all users: %s\n", err)
		} else {
			returnStr += " and users reset in DB"
		}
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(returnStr))
}