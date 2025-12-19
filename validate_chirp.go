package main

import (
	"net/http"
	"encoding/json"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type requestParams struct {
		Body string `json:"body"`
	}

	type jsonResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	w.Header().Set("Content-Type", "application/json")

	const chirpMaxLen = 140

	decoder := json.NewDecoder(req.Body)
	reqBody := requestParams{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}

	if len(reqBody.Body) > chirpMaxLen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, jsonResponse{
		CleanedBody: cleanChirp(reqBody.Body),
	})
}

func cleanChirp(chirp string) string {
	profaneList := [3]string{"kerfuffle", "sharbert", "fornax"}
	
	words := strings.Split(chirp, " ")

	for i, word := range words {
		for _, profane := range profaneList {
			if strings.Contains(strings.ToLower(word), profane) {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}
