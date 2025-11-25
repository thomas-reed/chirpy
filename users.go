package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/mail"
	"time"

	"github.com/google/uuid"
)

type user struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) addUserHandler(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	user := user{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding user parameters", err)
		return
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, "Email is not valid", err)
		return
	}
	ctx := context.Background()

	userFromDB, err := cfg.db.CreateUser(ctx, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding user to db", err)
		return
	}
	user.Id = userFromDB.ID
	user.CreatedAt = userFromDB.CreatedAt
	user.UpdatedAt = userFromDB.UpdatedAt
	user.Email = userFromDB.Email
	jsonResponse(w, http.StatusCreated, user)
}