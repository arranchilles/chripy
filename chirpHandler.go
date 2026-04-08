package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

type ChirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	Userid    uuid.UUID `json:"user_id"`
}

type chirpsResponse []ChirpResponse

func (c chirpsResponse) sortByCreatedAt(direction string) {
	sort.Slice(c, func(i, j int) bool {
		if direction == "desc" {
			return c[i].CreatedAt.After(c[j].CreatedAt)
		} else {
			return c[i].CreatedAt.Before(c[j].CreatedAt)
		}
	})
}

func (c *apiConfig) chirpPostHandler(writer http.ResponseWriter, request *http.Request) {
	chirpData, err := jsonDecode[chirp](request.Body)
	if err != nil {
		fmt.Print(err)
	}

	token, err := auth.GetBearerToken(request.Header)

	if err != nil {
		respondWithError(writer, 400, "Invalid token")
	}

	user, err := auth.ValidateJWT(token, c.secret)

	if err != nil {
		respondWithError(writer, 401, err.Error())
		return
	}

	err = validate_chirp(chirpData)

	if err != nil {
		respondWithError(writer, 400, err.Error())
		return
	}

	newBody := censorText(chirpData.Body)

	args := database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   newBody,
		UserID: user,
	}
	newChirp, err := c.dbQueries.CreateChirp(context.Background(), args)

	if err != nil {
		respondWithError(writer, 500, err.Error())
		return
	}

	responseData := ChirpResponse{}
	responseData.ID = newChirp.ID
	responseData.CreatedAt = newChirp.CreatedAt
	responseData.UpdatedAt = newChirp.UpdatedAt
	responseData.Body = newChirp.Body
	responseData.Userid = newChirp.UserID
	respondWithJSON(writer, 201, responseData)
}

func (c *apiConfig) chirpGetHandler(writer http.ResponseWriter, r *http.Request) {
	var responseData chirpsResponse
	var chirps []database.Chirp
	var err error

	if id := r.URL.Query().Get("author_id"); id != "" {
		uuid, err := uuid.Parse(id)
		if err != nil {
			errorHandler(writer, err)
		}
		chirps, err = c.dbQueries.GetChirpsByAuthor(context.Background(), uuid)
	} else {
		chirps, err = c.dbQueries.GetChirps(context.Background())
	}

	if err != nil {
		fmt.Print(err.Error())
		respondWithError(writer, 500, "Server Error")
		return
	}

	for _, chirpData := range chirps {
		item := ChirpResponse{}
		item.ID = chirpData.ID
		item.CreatedAt = chirpData.CreatedAt
		item.UpdatedAt = chirpData.UpdatedAt
		item.Body = chirpData.Body
		item.Userid = chirpData.UserID
		responseData = append(responseData, item)
	}

	if sort := r.URL.Query().Get("sort"); sort != "" {
		responseData.sortByCreatedAt(sort)
	}

	respondWithJSON(writer, 200, responseData)
}

func (c *apiConfig) chirpGetSingleHandler(writer http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")

	if id == "" {
		respondWithError(writer, 404, "No Chirp ID")
		return
	}
	chirpId, err := uuid.Parse(id)

	if err != nil {
		respondWithError(writer, 404, "Invalid ID format")
		return
	}

	chirp, err := c.dbQueries.GetChirp(context.Background(), chirpId)

	if err != nil {
		respondWithError(writer, 404, fmt.Sprintf("No chirp with id of: %s", id))
		fmt.Print(err.Error())
		return
	}

	responseData := ChirpResponse{
		ID:        chirpId,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		Userid:    chirp.UserID,
	}

	respondWithJSON(writer, 200, responseData)

}

func (c *apiConfig) chirpDeleteSingleHandler(writer http.ResponseWriter, request *http.Request) {

	userID, err := c.checkUserIdentity(request)

	if err != nil {
		errorHandler(writer, NewAppError(ErrorTypeAuth, "Unorthorized Access to Endpoint", err))
		return
	}

	chirpURLValue := request.PathValue("chirpID")

	if chirpURLValue == "" {
		errorHandler(writer, NewAppError(ErrorTypeValidation, "Need Chirp ID in the last url field", err))
		return
	}

	chirpID, err := uuid.Parse(chirpURLValue)

	if err != nil {
		errorHandler(writer, NewAppError(ErrorTypeValidation, "Invlaid chirp ID", err))
		return
	}

	chirp, err := c.dbQueries.GetChirp(context.Background(), chirpID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorHandler(writer, NewAppError(ErrorTypeResourceNotFound, "No Chirp by that ID", err))
			return
		}
		errorHandler(writer, err)
		return
	}

	if chirp.UserID != userID {
		errorHandler(
			writer,
			NewAppError(
				ErrorTypeForbidden,
				"Requested Chirp does not belong to user",
				fmt.Errorf("requested chrip belongs to %v. Unortherized by %v", chirp.UserID.String(), userID.String()),
			),
		)
		return
	}

	err = c.dbQueries.DeleteChirp(context.Background(), chirpID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorHandler(writer, NewAppError(ErrorTypeResourceNotFound, "No Chirp by that ID", err))
			return
		}
		errorHandler(writer, err)
		return
	}

	respondWithJSON(writer, 204, nil)
}
