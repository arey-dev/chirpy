package auth

import (
	"github.com/alexedwards/argon2id"
	"fmt"
	"errors"
	"net/http"
	"strings"
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)

	if err != nil {
		return "", fmt.Errorf("Error hashing password: %v", err)
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error comparing password: %v", err)
	}

	return match, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("No Authorization Header Found")
	}

	authHeader = strings.TrimSpace(authHeader)

	headerSlice := strings.Split(authHeader, " ")

	// check length
	if len(headerSlice) != 2 {
		return "", errors.New("Invalid API Key")
	}

	// get token
	apiKey := headerSlice[1]

	// check if bearer or valid
	if headerSlice[0] != "ApiKey" || apiKey == "" {
		return "", errors.New("No API Key found")
	}

	return apiKey, nil
}
