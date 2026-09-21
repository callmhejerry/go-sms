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
	Details    any    `json:"details"`
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
func New(code, message string, status int, err error, details any) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Err:        err,
		Details:    details,
	}
}

func Wrap(err error, code, message string, status int, details any) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Err:        err,
		Details:    details,
	}
}

const (
	CheckViolation      = "check_violation"
	UniqueViolation     = "unique_violation"
	ForeignKeyViolation = "foreign_key_violation"
	NotNullViolation    = "not_null_violation"
)

var (
	ErrNotFound = New("not_found", "resource not found", http.StatusNotFound, nil, nil)

	ErrConflict = New("conflict", "resource already exists", http.StatusConflict, nil, nil)

	ErrValidation = New("validation_error", "validation failed", http.StatusBadRequest, nil, nil)

	ErrUnauthorized = New("unauthorized", "authentication required", http.StatusUnauthorized, nil, nil)

	ErrForbidden = New("forbidden", "you do not have permission to perform this action", http.StatusForbidden, nil, nil)

	ErrInternal = New("internal_error", "internal server error", http.StatusInternalServerError, nil, nil)
)

// Convenience helpers
func NotFound(message string) *AppError {
	return New("not_found", message, http.StatusNotFound, nil, nil)
}

func Conflict(message string) *AppError {
	return New("conflict", message, http.StatusConflict, nil, nil)
}

func Validation(message string) *AppError {
	return New("validation_error", message, http.StatusBadRequest, nil, nil)
}

func Internal(err error, message string) *AppError {
	return Wrap(err, "internal_error", message, http.StatusInternalServerError, nil)
}
