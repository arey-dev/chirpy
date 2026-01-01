package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
	"fmt"
	"net/http"
)

func TestCreateValidateJWT(t *testing.T) {
	mySigningKey := "mySigningKey"
	newID := uuid.New()
	expiresIn, _ := time.ParseDuration("1h")

	token, err := MakeJWT(newID, mySigningKey, expiresIn)

	if err != nil {
		t.Errorf("Error Making JWT: %v", err)
	}

	_, err = ValidateJWT(token, mySigningKey)

	if err != nil {
		t.Errorf("Error validating JWT: %v", err)
	}
}

func TestGetBearerToken(t *testing.T) {
	headers := make(http.Header)
	token := "token"
	headers.Set("Authorization", fmt.Sprintf("  Bearer %s  ", token))

	headers2 := make(http.Header)
	headers2.Set("Content-Type", "application/json")

	headers3 := make(http.Header)
	token3 := ""
	headers3.Set("Authorization", fmt.Sprintf("Bearer %s", token3))

	headers4 := make(http.Header)
	headers4.Set("Authorization", "apikey")

	headers5 := make(http.Header)
	headers5.Set("Authorization", "Bearer two three")

	cases := []struct{
		input http.Header
		name string
		expected string
		wantErr bool
	}{
		{
			input: headers,
			name: "returns bearer token",
			expected: token,
			wantErr: false,
		},
		{
			input: headers2,
			name: "no authorization",
			expected: "",
			wantErr: true,
		},
		{
			input: headers3,
			name: "no token",
			expected: "",
			wantErr: true,
		},
		{
			input: headers4,
			name: "bearer token not set",
			expected: "",
			wantErr: true,
		},
		{
			input: headers5,
			name: "invalid bearer token",
			expected: "",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual, err := GetBearerToken(c.input)

			if c.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if actual != c.expected {
				t.Errorf("expected %s, actual %s", c.expected, actual)
			}
		})
	}
}
