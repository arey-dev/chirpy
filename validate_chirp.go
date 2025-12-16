package main

import (
	"net/http"
	"encoding/json"
	"log"
)

type ValidateChirpyRequest struct {
	Body string `json:"body"`
}

type APIResponse struct {
	Valid bool `json:"valid"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func handlerValidateChirp(w http.ResponseWriter, req *http.Request) {	
	w.Header().Set("Content-Type", "application/json")

	chirpMaxLen := 140

	decoder := json.NewDecoder(req.Body)
	reqBody := ValidateChirpyRequest{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding request body: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Something went wrong",
		})
		return
	}

	if len(reqBody.Body) > chirpMaxLen {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Chirp is too long",
		})
		return
	} 

	data, err := json.Marshal(APIResponse{
		Valid: true,
	})
	if err != nil {
		log.Printf("Error marshalling api response: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: "Something went wrong",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
