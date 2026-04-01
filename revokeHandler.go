package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func (c *apiConfig) revokeHandler(writer http.ResponseWriter, request *http.Request) {
	token, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(writer, 401, "Invalid Token")
		fmt.Println(err.Error())
		return
	}

	params := database.RevokeTokenParams{
		Token: token,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err = c.dbQueries.RevokeToken(context.Background(), params)

	if err != nil {
		respondWithError(writer, 401, "Invalid Token")
		fmt.Println(err.Error())
		return
	}

	writer.WriteHeader(204)
}
