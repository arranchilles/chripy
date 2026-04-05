package main

import (
	"chirpy/internal/auth"
	"context"
	"fmt"
	"net/http"
	"time"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (c *apiConfig) refreshHandler(writer http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(writer, 500, "Server issue")
		fmt.Println(err.Error())
		return
	}

	user, err := c.dbQueries.GetUserFromRefreshToken(context.Background(), refreshToken)

	if user.RevokedAt.Valid {
		err = fmt.Errorf("Token already Revoked")
	}

	if err != nil {
		respondWithError(writer, 401, "Invalid Token")
		fmt.Println(err.Error())
		return
	}

	accessToken, err := auth.MakeJWT(user.UserID, c.secret, time.Hour)

	if err != nil {
		respondWithError(writer, 500, "Server issue")
		fmt.Println(err.Error())
		return
	}

	newAccessToken := RefreshResponse{
		Token: accessToken,
	}

	respondWithJSON(writer, 200, newAccessToken)
}
