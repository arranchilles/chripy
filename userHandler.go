package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type userInfo struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ExpiresIn int64  `json:"expires_in_seconds"`
}

func (c *apiConfig) handleUsersPost(w http.ResponseWriter, r *http.Request) {

	postData, err := jsonDecode[userInfo](r.Body)

	if err != nil {
		respondWithError(w, 500, "Server Could not Process JSON")
		return
	}

	if postData.Password == "" {
		respondWithError(w, 400, "No Password sent")
		return
	}

	hash, err := auth.HashPassword(postData.Password)

	if err != nil {
		respondWithError(w, 500, "Server Error")
		return
	}

	userInfo := database.CreateUserParams{
		Email:    postData.Email,
		Password: hash,
	}

	user, err := c.dbQueries.CreateUser(context.Background(), userInfo)

	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	responseData := UserResponse{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, 201, responseData)
}

func (c *apiConfig) handleUsersPut(writer http.ResponseWriter, request *http.Request) {

	userID, err := c.checkUserIdentity(request)

	if err != nil {
		errorHandler(writer, err)
		return
	}

	newInfo, err := jsonDecode[userInfo](request.Body)

	if err != nil {
		errorHandler(writer, err)
		return
	}

	hash, err := auth.HashPassword(newInfo.Password)

	if err != nil {
		errorHandler(writer, err)
		return
	}

	args := database.UpdateUserInfoParams{
		ID:       userID,
		Email:    newInfo.Email,
		Password: hash,
	}

	user, err := c.dbQueries.UpdateUserInfo(context.Background(), args)

	if err != nil {
		errorHandler(writer, NewAppError(ErrorTypeAuth, "Invalid Operation", err))
		return
	}

	responseData := UserResponse{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(writer, 200, responseData)
}
