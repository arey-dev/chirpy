package main

import (
	"net/http"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"github.com/arey-dev/chirpy/internal/database"
	"github.com/arey-dev/chirpy/internal/auth"
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
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	reqBody := requestBody{}
	err := decoder.Decode(&reqBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	createUserParams := database.CreateUserParams{
		Email: reqBody.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(req.Context(),createUserParams)

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

