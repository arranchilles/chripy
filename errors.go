package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeAuth       ErrorType = "auth"
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeInternal   ErrorType = "internal"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (a *AppError) Error() string {
	return a.Message
}

func (a *AppError) Log() {
	if a.Err != nil {
		log.Println(a.Err)
		return
	}
	log.Println(a.Message)
}

func NewAppError(errorType ErrorType, message string, err error) error {
	return &AppError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

func errorHandler(writer http.ResponseWriter, err error) {
	var appError *AppError
	if errors.As(err, &appError) {
		appErrorHandler(writer, appError)
		return
	}
	fmt.Print(err.Error())
	respondWithError(writer, 500, "Internal Server Error")
}

func appErrorHandler(writer http.ResponseWriter, err *AppError) {
	err.Log()
	switch err.Type {
	case ErrorTypeAuth:
		respondWithError(writer, 401, err.Message)
	case ErrorTypeValidation:
		respondWithError(writer, 400, err.Message)
	case ErrorTypeInternal:
		respondWithError(writer, 500, err.Message)
	default:
		respondWithError(writer, 500, "Internal Server Error")
	}
}
