package main

import (
	"strings"
	"errors"
)

func validateChirp(chirp string) (string, error) {
	const chirpMaxLen = 140

	if len(chirp) > chirpMaxLen {
		return "", errors.New("Chirp is too long")
	}

	return cleanChirp(chirp), nil
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
