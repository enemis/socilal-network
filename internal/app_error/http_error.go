package app_error

import "net/http"

type AppError struct {
	message       string
	field         string
	status        int
	originalError error
}

const SERVER_ERROR_MESSAGE = "Something went wrong"

func (e AppError) Error() string {
	return e.message
}

func (e AppError) Status() int {
	return e.status
}

func (e AppError) OriginalError() error {
	return e.originalError
}

func New(originalError error, message string, status int) *AppError {
	if message == "" {
		message = SERVER_ERROR_MESSAGE
	}

	if status == 0 {
		status = http.StatusInternalServerError
	}

	return &AppError{message: message, status: status, originalError: originalError}
}

func NewBadRequestFromError(originalError error) *AppError {
	return NewFromError(originalError, http.StatusBadRequest)
}

func NewForbiddenFromError(originalError error) *AppError {
	return NewFromError(originalError, http.StatusForbidden)
}

func NewNotFoundFromError(originalError error) *AppError {
	return NewFromError(originalError, http.StatusNotFound)
}

func NewFromError(originalError error, status int) *AppError {
	if status == 0 {
		status = http.StatusInternalServerError
	}

	return &AppError{message: originalError.Error(), status: status, originalError: originalError}
}

func NewInternalServerError(originalError error) *AppError {
	return New(originalError, "", 0)
}
