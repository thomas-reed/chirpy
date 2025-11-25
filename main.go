package main

import (
	"os"
	"log"
	"net/http"
	"sync/atomic"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/thomas-reed/chirpy/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not load env file")
	}
	
	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Fatalln("Port not found")
	}
	filepathRoot := os.Getenv("FILEPATHROOT")
	if filepathRoot == "" {
		log.Fatalln("filepath root not found")
	}
	platform := os.Getenv("PLATFORM")
	
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("Database connection string not found")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Could not open database connection")
	}

	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db: database.New(dbConn),
		platform: platform,
	}

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(fsHandler))
	mux.HandleFunc("GET /admin/metrics", apiConfig.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetHandler)
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	mux.HandleFunc("POST /api/users", apiConfig.addUserHandler)

	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Server listening on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}