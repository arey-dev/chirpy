package auth

import (
	 "github.com/golang-jwt/jwt/v5"
	 "github.com/google/uuid"
	 "time"
	 "fmt"
	 "net/http"
	 "strings"
	 "errors"
	 "encoding/hex"
	 "crypto/rand"
)

type ChirpyClaims struct {
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := ChirpyClaims{
		jwt.RegisteredClaims{
			Issuer: "chirpy",
			Subject: userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", fmt.Errorf("Error signing token: %v", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ChirpyClaims{}, func (token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error parsing token: %v\n", err)
	}

	claims := token.Claims.(*ChirpyClaims)

	subject, err := claims.GetSubject()

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error parsing jwt subject: %v\n", err)
	}

	userID, err := uuid.Parse(subject)

	if err != nil {
		return uuid.Nil, fmt.Errorf("Error parsing uuid: %v\n", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("No Authorization Header Found")
	}

	authHeader = strings.TrimSpace(authHeader)

	headerSlice := strings.Split(authHeader, " ")

	// check length
	if len(headerSlice) != 2 {
		return "", errors.New("Invalid Bearer Token")
	}

	// get token
	token := headerSlice[1]

	// check if bearer or valid
	if headerSlice[0] != "Bearer" || token == "" {
		return "", errors.New("No Bearer Token found")
	}

	return token, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	encodedKey := hex.EncodeToString(key)
	return encodedKey, nil
}
