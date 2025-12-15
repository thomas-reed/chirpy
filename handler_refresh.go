package main

import (
	"net/http"
	"time"

	"github.com/thomas-reed/chirpy/internal/auth"
)

type refresh struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request){
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Refresh token not in headers", err)
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found in DB", err)
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) || refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token no longer valid", err)
		return
	}
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating access token", err)
		return
	}
	jsonResponse(w, http.StatusOK, refresh{
		Token: accessToken,
	})
}