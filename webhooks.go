package main

import (
	"chirpy/internal/auth"
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type WebHook struct {
	Event string      `json:"event"`
	Data  webHookData `json:"data"`
}

type webHookData struct {
	UserID uuid.UUID `json:"user_id"`
}

func (c *apiConfig) webhooksPostHandler(writer http.ResponseWriter, request *http.Request) {

	apiKey, err := auth.GetAPIKey(request.Header)

	if err != nil {
		errorHandler(writer, NewAppError(ErrorTypeAuth, err.Error(), err))
		return
	}

	if apiKey != c.polkaKey {
		errorHandler(writer, NewAppError(ErrorTypeAuth, "Invalid API KEY", nil))
		return
	}

	webHook, err := jsonDecode[WebHook](request.Body)
	if err != nil {
		errorHandler(writer, NewAppError(ErrorTypeValidation, "Webhook Post JSON is Invalid", err))
	}

	if webHook.Event != "user.upgraded" {
		respondWithJSON(writer, 204, nil)
		return
	}

	err = c.dbQueries.AddUserToChirpRed(request.Context(), webHook.Data.UserID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorHandler(writer, NewAppError(ErrorTypeResourceNotFound, "There is no user by that ID", err))
			return
		}
		errorHandler(writer, err)
		return
	}

	respondWithJSON(writer, 204, nil)
}
