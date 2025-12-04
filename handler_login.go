package main

import (
	"time"
	"encoding/json"
	"net/http"
	"github.com/thomas-reed/chirpy/internal/auth"
)

const maxExpireTime = int64(time.Hour)

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int64 `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding user parameters", err)
		return
	}
	
	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching user data", err)
		return
	}
	success, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !success || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	expires := params.ExpiresInSeconds
	if expires <= 0 || expires > maxExpireTime {
		expires = maxExpireTime
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(expires))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating JWT", err)
		return
	}
	
	jsonResponse(w, http.StatusOK, User{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	})
}