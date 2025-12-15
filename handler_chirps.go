package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thomas-reed/chirpy/internal/auth"
	"github.com/thomas-reed/chirpy/internal/database"
)

type Chirp struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) addChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding chirp parameters", err)
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting authorization token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
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

	newChirpParams := database.CreateChirpParams{
		Body: cleanText(params.Body),
		UserID: userID,
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), newChirpParams)
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


func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	rawChirps, err := cfg.db.GetAllChirps(r.Context())
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

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), id)
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

func (cfg *apiConfig) deleteChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not in headers", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}
	idStr := r.PathValue("chirpID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id", err)
		return
	}
	getChirpParams := database.GetChirpByIDAndUserParams{
		ID: id,
		UserID: userID,
	}
	_, err = cfg.db.GetChirpByIDAndUser(r.Context(), getChirpParams)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Must be author of chirp to delete", err)
		return
	}

	deleteParams := database.DeleteChirpByIdAndUserParams{
		ID: id,
		UserID: userID,
	}
	_, err = cfg.db.DeleteChirpByIdAndUser(r.Context(), deleteParams)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "No chirp with provided id", err)
		return
	}

	jsonResponse(w, http.StatusNoContent, struct{}{})
}