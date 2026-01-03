package main

import (
	"net/http"
	"encoding/json"
	"strings"
	"github.com/google/uuid"
	"github.com/arey-dev/chirpy/internal/auth"
)

type EventType string

const EventUserUpgraded EventType = "user.upgraded"

func handleWebhooks(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", err)
		return
	}

	type RequestBody struct {
		Event string `json:"event"`
		Data struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(req.Body)
	reqBody := RequestBody{}
	err = decoder.Decode(&reqBody)
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
