package main

import (
	"net/http"
	"encoding/json"
)

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {	
	type requestParams struct {
		Body string `json:"body"`
	}

	type jsonResponse struct {
		Valid bool `json:"valid"`
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
		Valid: true,
	})
}
