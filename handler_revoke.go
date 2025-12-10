package main

import (
	"net/http"

	"github.com/thomas-reed/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, req *http.Request){
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Refresh token not in headers", err)
		return
	}
	_, err = cfg.db.RevokeRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Token not found in DB", err)
		return
	}
	
	jsonResponse(w, http.StatusNoContent, struct{}{})
}