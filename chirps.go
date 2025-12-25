package main

import (
	"net/http"
	"github.com/google/uuid"
	"encoding/json"
	"time"
	"github.com/arey-dev/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func createChirp(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Params{}
	err := decoder.Decode(&params)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	validChirp, err := validateChirp(params.Body)

	if err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, err.Error(), err)
		return
	}

	params.Body = validChirp
	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body: validChirp,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Creating Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func getAllChirps(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching chirps", err)
		return
	}

	chirpsRes := []Chirp{}

	for _, chirp := range chirps {
		chirpsRes = append(chirpsRes, Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpsRes)
}
