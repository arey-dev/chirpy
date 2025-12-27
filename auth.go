package main

import (
	"net/http"
	"encoding/json"
	"github.com/arey-dev/chirpy/internal/auth"
)

func loginUser(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Params{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	isPasswordMatch, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)

	if err != nil || !isPasswordMatch {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}
	
	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email, 
	})
}
