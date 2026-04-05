package main

import (
	"chirpy/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type chirp struct {
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated-at"`
}

type Validated_chirp struct {
	CleanedBody string `json:"cleaned_body"`
}

type ErrorPayload struct {
	Error string `json:"error"`
}

const metricsTemplate = "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>"

func (cfg *apiConfig) handleMetrics(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	body := []byte(fmt.Sprintf(metricsTemplate, cfg.fileserverHits.Load()))
	writer.Write(body)
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.Header().Add("Cache-Control", "no-cache")
}

func (cfg *apiConfig) handleReset(writer http.ResponseWriter, request *http.Request) {

	if cfg.platfrom != "dev" {
		respondWithError(writer, 403, "Invalid Permission")
		return
	}

	writer.WriteHeader(200)
	body := []byte("Reset application data")
	cfg.fileserverHits = atomic.Int32{}
	cfg.dbQueries.DeleteUsers(context.Background())
	writer.Write(body)
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.Header().Add("Cache-Control", "no-cache")
}

func handleHealthz(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	body := []byte("200 OK")
	writer.Write(body)
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.Header().Add("Cache-Control", "no-cache")
}

func handleChirpValidate(writer http.ResponseWriter, request *http.Request) {
	chirpData, err := jsonDecode[chirp](request.Body)
	if err != nil {
		fmt.Print(err)
	}

	err = validate_chirp(chirpData)

	if err != nil {
		respondWithError(writer, 400, err.Error())
		return
	}

	newBody := censorText(chirpData.Body)

	responseData := Validated_chirp{
		CleanedBody: newBody,
	}

	respondWithJSON(writer, 200, responseData)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "text/json")
	errorPayload := ErrorPayload{
		Error: msg,
	}
	messageBody, err := json.Marshal(errorPayload)
	if err != nil {
		fmt.Print(err)
	}
	w.Write(messageBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "text/json")
	responsePayload, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, 500, "Failed to process Chirp")
		fmt.Print(err, "respondjson error")
	}
	w.Write(responsePayload)
}

func (c *apiConfig) checkUserIdentity(request *http.Request) (uuid.UUID, error) {

	token, err := auth.GetBearerToken(request.Header)

	if err != nil {
		return uuid.UUID{}, &AppError{
			Type:    ErrorTypeAuth,
			Message: "No Authorization Token in header",
			Err:     err,
		}
	}

	userID, err := auth.ValidateJWT(token, c.secret)

	if err != nil {
		return uuid.UUID{}, &AppError{
			Type:    ErrorTypeAuth,
			Message: "Authorization token in header is Invalid",
			Err:     err,
		}
	}
	return userID, nil
}
