package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"errors"
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
		return
	}

	if userData.Email == "" {
		respondWithError(writer, 400, "Email needed in JSON body")
		return
	}

	if userData.Password == "" {
		respondWithError(writer, 400, "Password needed in JSON body")
		return
	}

	expiresIn = time.Duration(userData.ExpiresIn) * time.Second

	if expiresIn.Seconds() == 0 {
		expiresIn = time.Hour
	}

	if expiresIn.Hours() > 1 {
		expiresIn = time.Hour
	}

	fmt.Println(expiresIn)

	user, err := c.validateUser(userData)

	if err != nil {
		errorHandler(writer, err)
		return
	}

	token, err := auth.MakeJWT(user.ID, c.secret, time.Duration(expiresIn))

	if err != nil {
		errorHandler(writer, NewAppError(ErrorTypeInternal, "Server Error", err))
		return
	}

	refreshToken := auth.MakeRefreshToken()

	err = c.addRefeshToken(refreshToken, user.ID)

	if err != nil {
		respondWithError(writer, 500, "Server Error")
		return
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

func (c *apiConfig) validateUser(userInput userInfo) (database.User, error) {

	user, err := c.dbQueries.GetUser(context.Background(), userInput.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.User{}, NewAppError(ErrorTypeValidation, "Invalid Email or Password", err)
		}

		return database.User{}, NewAppError(ErrorTypeInternal, "Internal Server Error", err)

	}

	isPassword, err := auth.CheckPasswordHash(userInput.Password, user.Password)

	if err != nil {
		return database.User{}, NewAppError(ErrorTypeInternal, "Internal Server Error", err)
	}

	if isPassword != true {
		return database.User{}, NewAppError(ErrorTypeValidation, "Invalid Email or Password", nil)
	}

	return user, nil
}
