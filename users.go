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

func createUser(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
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

func updateUser(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	UserID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}
	
	type requestBody struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	reqBody := requestBody{}
	err = decoder.Decode(&reqBody)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}
	
	hashedPassword, err := auth.HashPassword(reqBody.Password)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	updateUserParams := database.UpdateUserParams{
		Email: reqBody.Email,
		HashedPassword: hashedPassword,
		ID: UserID,
	}

	user, err := cfg.db.UpdateUser(req.Context(), updateUserParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Updating User", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email, 
	})
}
