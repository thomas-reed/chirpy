package main

import (
	"net/http"
	"encoding/json"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

type cleanedChirp struct {
	Cleaned string `json:"cleaned_body"`
}

func validateChirpHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	chirpToValidate := chirp{}
	err := decoder.Decode(&chirpToValidate)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding chirp parameters", err)
		return
	}
	if len(chirpToValidate.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	cleanedText := cleanText(chirpToValidate.Body)
	cleanedJson := cleanedChirp{
		Cleaned: cleanedText,
	}
	jsonResponse(w, http.StatusOK, cleanedJson)
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
