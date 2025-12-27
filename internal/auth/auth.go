package auth

import (
	"github.com/alexedwards/argon2id"
	"fmt"
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
