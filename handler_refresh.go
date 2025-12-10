package main

import (
	"net/http"
	"time"

	"github.com/thomas-reed/chirpy/internal/auth"
)

type refresh struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request){
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Refresh token not in headers", err)
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), token)
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