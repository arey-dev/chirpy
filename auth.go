package main

import (
	"net/http"
	"encoding/json"
	"strconv"
	"database/sql"
	"time"
	"errors"
	"github.com/arey-dev/chirpy/internal/auth"
	"github.com/arey-dev/chirpy/internal/database"
)

type UserWithToken struct {
	User
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

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

	tokenExpiresIn, err := setExpiresIn(cfg.jwtTTL)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, tokenExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	refreshToken, _ := auth.MakeRefreshToken()
	createRefreshTokenParams := database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	_, err = cfg.db.CreateRefreshToken(req.Context(), createRefreshTokenParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error Creating Refresh Token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, UserWithToken{
		User: User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email, 
		},
		Token: token,
		RefreshToken: refreshToken,
	})
}

func issueNewToken(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	// check if user is authenticated
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	refreshToken, err := cfg.db.GetUserFromRefreshToken(req.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error generating refresh token", err)
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Error generating refresh token", err)
		return
	}

	tokenExpiresIn, err := setExpiresIn(cfg.jwtTTL)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error generating token", err)
		return
	}

	newToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, tokenExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct{
		Token string `json:"token"`
	}{
		Token: newToken,
	})
}

func revokeToken(cfg *apiConfig, w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
		return
	}

	params := database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		Token: token,
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error revoking refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}

func setExpiresIn (defaultTTL string) (time.Duration, error) { 
	defaultExpiresIn, err := strconv.Atoi(defaultTTL)
	if err != nil {
		return 0, errors.New("error converting jwt ttl string to int")
	}
	return time.Duration(defaultExpiresIn)*time.Second, nil
}
