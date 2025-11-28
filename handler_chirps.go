package main

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"
	"github.com/google/uuid"
	"github.com/thomas-reed/chirpy/internal/database"
)

type Chirp struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) addChirpHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserId string `json:"user_id"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding chirp parameters", err)
		return
	}
	if len(params.Body) == 0  {
		respondWithError(w, http.StatusBadRequest, "Chirp is empty", nil)
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	id, err := uuid.Parse(params.UserId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user id", err)
		return
	}

	newChirpParams := database.CreateChirpParams{
		Body: cleanText(params.Body),
		UserID: id,

	}

	chirp, err := cfg.db.CreateChirp(req.Context(), newChirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp", err)
		return
	}
	
	jsonResponse(w, http.StatusCreated, Chirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	})
}

func cleanText(text string) string {
	badWords := map[string]struct{} {
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	replacementText := "****"
	words := strings.Split(text, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = replacementText
		}
	}
	return strings.Join(words, " ")
}


func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, req *http.Request) {
	rawChirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting all chirps", err)
		return
	}

	chirps := make([]Chirp, 0, len(rawChirps))
	for _, chirp := range rawChirps {
		chirps = append(chirps, Chirp{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		})
	}
	
	jsonResponse(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, req *http.Request) {
	idStr := req.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(req.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp with provided id", err)
		return
	}

	jsonResponse(w, http.StatusOK, Chirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	})
}