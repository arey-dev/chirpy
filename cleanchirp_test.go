package main

import (
	"fmt"
	"testing"
)

func TestCleanChirp(t *testing.T) {
	cases := []struct{
		input string
		expected string
	}{
		{
			input: "",
			expected: "",
		},
		{
			input:  "This is a kerfuffle opinion I need to share with the world",
			expected: "This is a **** opinion I need to share with the world",
		},
		{
			input: "I really need a kerfuffle to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
		{
			input: "I really need a kerfufflefornax to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i + 1), func(t *testing.T) {
			actual := cleanChirp(c.input)

			if c.input == "" && actual != "" {
				t.Error("must return empty if the input is empty")
				return
			}

			if actual != c.expected {
				t.Errorf("case %d:\n input: %s\n expected: %s\n actual: %s\n", i, c.input, c.expected, actual)
			}
		})
	}
}
