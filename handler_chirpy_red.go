package main

import (
	"net/http"
	"encoding/json"
	"errors"
	"database/sql"

	"github.com/google/uuid"
	"github.com/thomas-reed/chirpy/internal/auth"
)

func (cfg *apiConfig) addChirpyRedHandler(w http.ResponseWriter, r *http.Request){
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiToken, err := auth.GetAPIKey(r.Header)
	if err != nil || apiToken != cfg.polkaKey  {
		respondWithError(w, http.StatusUnauthorized, "Error getting API token", err)
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding user parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		jsonResponse(w, http.StatusNoContent, struct{}{})
		return
	}
	_, err = cfg.db.AddChirpyRedByUserID(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
      jsonResponse(w, http.StatusNotFound, struct{}{})
    } else {
      jsonResponse(w, http.StatusInternalServerError, struct{}{})
    }
		return
	}
	
	jsonResponse(w, http.StatusNoContent, struct{}{})
}