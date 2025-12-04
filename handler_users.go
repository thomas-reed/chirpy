package main

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"time"
	"github.com/google/uuid"
	"github.com/thomas-reed/chirpy/internal/auth"
	"github.com/thomas-reed/chirpy/internal/database"
)

type User struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) addUserHandler(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding user parameters", err)
		return
	}
	if _, err := mail.ParseAddress(params.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, "Email is not valid", err)
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password cannot be empty", nil)
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	createUserParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	}

	user, err := cfg.db.CreateUser(req.Context(), createUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}
	
	jsonResponse(w, http.StatusCreated, User{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	})
}