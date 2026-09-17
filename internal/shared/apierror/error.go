package apierror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Helper constructors
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

func Wrap(err error, code, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Err:        err,
	}
}

var (
	ErrNotFound = New("not_found", "resource not found", http.StatusNotFound)

	ErrConflict = New("conflict", "resource already exists", http.StatusConflict)

	ErrValidation = New("validation_error", "validation failed", http.StatusBadRequest)

	ErrUnauthorized = New("unauthorized", "authentication required", http.StatusUnauthorized)

	ErrForbidden = New("forbidden", "you do not have permission to perform this action", http.StatusForbidden)

	ErrInternal = New("internal_error", "internal server error", http.StatusInternalServerError)
)

// Convenience helpers
func NotFound(message string) *AppError {
	return New("not_found", message, http.StatusNotFound)
}

func Conflict(message string) *AppError {
	return New("conflict", message, http.StatusConflict)
}

func Validation(message string) *AppError {
	return New("validation_error", message, http.StatusBadRequest)
}

func Internal(err error, message string) *AppError {
	return Wrap(err, "internal_error", message, http.StatusInternalServerError)
}
