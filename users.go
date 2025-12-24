package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func handlerCreateUser(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type requestBody struct {
		Email string `json:email`
	}

	decoder := json.NewDecoder(req.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), reqBody.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Creating User", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email, 
	})
}

