package main

import (
	"os"
	"log"
	"net/http"
	"sync/atomic"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/thomas-reed/chirpy/internal/database"
	"github.com/pressly/goose/v3"
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
	if err := goose.SetDialect("postgres"); err != nil {
    log.Fatalf("Error running goose SetDialect: %v", err)
	}

	if err := goose.Up(dbConn, "sql/schema"); err != nil {
			log.Fatalf("Error running goose Up: %v", err)
	}

	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db: database.New(dbConn),
		platform: platform,
	}

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(fsHandler))
	
	mux.HandleFunc("GET /api/healthz", healthHandler)
	mux.HandleFunc("POST /api/users", apiConfig.addUserHandler)
	mux.HandleFunc("GET /api/chirps", apiConfig.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.getChirpByIDHandler)
	mux.HandleFunc("POST /api/chirps", apiConfig.addChirpHandler)

	mux.HandleFunc("GET /admin/metrics", apiConfig.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiConfig.resetHandler)

	s := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Server listening on port: %s\n", port)
	log.Fatal(s.ListenAndServe())
}