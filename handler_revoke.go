package main

import (
	"net/http"

	"github.com/thomas-reed/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Refresh token not in headers", err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token not found in DB", err)
		return
	}
	
	jsonResponse(w, http.StatusNoContent, struct{}{})
}