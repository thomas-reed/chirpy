package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type invalidChirp struct {
	Err string `json:"error"`
}

type cleanedChirp struct {
	Cleaned string `json:"cleaned_body"`
}

func respondWithError(w http.ResponseWriter, errCode int, errMsg string, err error) {
	if err != nil {
		log.Printf("%s: %v", errMsg, err)
	}
	invalidJson := invalidChirp{
		Err: errMsg,
	}
	jsonResponse(w, errCode, invalidJson)
}

func validResponse(w http.ResponseWriter, chirpText string) {
	cleanedText := cleanText(chirpText)
	cleanedJson := cleanedChirp{
		Cleaned: cleanedText,
	}
	jsonResponse(w, http.StatusOK, cleanedJson)
}

func jsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)
		return
	}
	w.WriteHeader(code)
	w.Write(json)
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
