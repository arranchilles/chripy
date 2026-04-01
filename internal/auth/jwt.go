package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	expires := now.Add(expiresIn)

	claim := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expires),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	JWT, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return JWT, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.Parse(id)
}

func GetBearerToken(headers http.Header) (string, error) {
	rawToken := headers.Get("Authorization")

	if rawToken == "" {
		return "", fmt.Errorf("No Authorization token")
	}

	if strings.HasPrefix(rawToken, "Bearer ") {
		return strings.TrimPrefix(rawToken, "Bearer "), nil
	}

	return rawToken, nil
}
