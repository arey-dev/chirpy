package main

import (
	"net/http"
	"github.com/google/uuid"
	"strings"
	"encoding/json"
	"time"
	"github.com/arey-dev/chirpy/internal/database"
	"github.com/arey-dev/chirpy/internal/auth"
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
	}

	decoder := json.NewDecoder(req.Body)
	params := Params{}
	err := decoder.Decode(&params)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	// check if user is authenticated
	bearerToken, err := auth.GetBearerToken(req.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	UserID, err := auth.ValidateJWT(bearerToken, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
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
		UserID: UserID,
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

func getChirp(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing path value", err)
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusNotFound, "Chirp Not Found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Error fetching chirps", err)
		return
	}


	respondWithJSON(w, http.StatusOK, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}

func deleteChirp(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	// check if user is authenticated
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

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing path value", err)
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusNotFound, "Chirp Not Found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, "Error deleting chirps", err)
		return
	}

	if chirp.UserID != UserID {
		respondWithError(w, http.StatusForbidden, "Error deleting chirps", err)
		return
	}

	err = cfg.db.DeleteChirp(req.Context(), chirpID)
	if chirp.UserID != UserID {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirps", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
