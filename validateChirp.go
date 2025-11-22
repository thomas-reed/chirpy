package main

import (
	"net/http"
	"encoding/json"
)

type chirp struct {
	Body string `json:"body"`
}

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	
	decoder := json.NewDecoder(req.Body)
	chirpToValidate := chirp{}
	err := decoder.Decode(&chirpToValidate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}
	if len(chirpToValidate.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	validResponse(w, chirpToValidate.Body)
}