package main

import (
	"fmt"
	"strings"
)

var censored = []string{"kerfuffle", "sharbert", "fornax"}

func validate_chirp(data chirp) error {

	if len(data.Body) > 140 {
		return fmt.Errorf("Chirp length is over 140 characters")
	}

	return nil
}

func censorText(text string) string {
	noCapText := strings.ToLower(text)
	words := strings.Split(noCapText, " ")
	originalWords := strings.Split(text, " ")
	for i, word := range words {
		for _, censor := range censored {
			if word == censor {
				originalWords[i] = "****"
			}
		}
	}
	censoredText := strings.Join(originalWords, " ")
	return censoredText
}
