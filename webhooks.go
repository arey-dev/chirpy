package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"github.com/google/uuid"
)

type EventType string

const EventUserUpgraded EventType = "user.upgraded"

func handleWebhooks(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	type RequestBody struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request", err)
		return
	}
	
	if reqBody.Event != string(EventUserUpgraded) {
		respondWithJSON(w, http.StatusNoContent, struct{}{})
		return
	}


	userID, err := uuid.Parse(reqBody.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing user id", err)
		return
	}

	_, err = cfg.db.UpdateUserToChirpyRed(req.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			respondWithError(w, http.StatusNotFound, "User Not Found", err)
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
