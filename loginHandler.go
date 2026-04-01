package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type LoginResponse struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (c *apiConfig) handleLogin(writer http.ResponseWriter, request *http.Request) {
	loginData := request.Body
	userData, err := jsonDecode[userInfo](loginData)
	var expiresIn time.Duration

	if err != nil {
		respondWithError(writer, 500, "Server Error")
	}

	if userData.Email == "" {
		respondWithError(writer, 400, "Email needed in JSON body")
	}

	if userData.Password == "" {
		respondWithError(writer, 400, "Password needed in JSON body")
	}

	expiresIn = time.Duration(userData.ExpiresIn) * time.Second

	if expiresIn.Seconds() == 0 {
		expiresIn = time.Hour
	}

	if expiresIn.Hours() > 1 {
		expiresIn = time.Hour
	}

	fmt.Println(expiresIn)
	user, err := c.dbQueries.GetUser(context.Background(), userData.Email)

	if err != nil {
		respondWithError(writer, 500, "Server Error")
		return
	}

	isPassword, err := auth.CheckPasswordHash(userData.Password, user.Password)

	if err != nil {
		respondWithError(writer, 500, "Server Error")
		return
	}

	if isPassword != true {
		respondWithError(writer, 401, "Incorrect Password")
		return
	}

	token, err := auth.MakeJWT(user.ID, c.secret, time.Duration(expiresIn))

	if err != nil {
		respondWithError(writer, 500, "Server Error")
	}

	refreshToken := auth.MakeRefreshToken()

	err = c.addRefeshToken(refreshToken, user.ID)

	if err != nil {
		respondWithError(writer, 500, "Server Error")
	}

	userResponse := LoginResponse{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(writer, 200, userResponse)
}

func (c *apiConfig) addRefeshToken(token string, userId uuid.UUID) error {
	createdAt := time.Now()
	expiresAt := createdAt.Add(time.Hour * 24 * 60)

	params := database.AddRefreshTokenParams{
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
		UserID:    userId,
		Token:     token,
	}

	err := c.dbQueries.AddRefreshToken(context.Background(), params)

	if err != nil {
		return err
	}

	return nil
}
