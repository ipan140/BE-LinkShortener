package errors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrUnauthorized   = errors.New("unauthorized access")
	ErrForbidden      = errors.New("forbidden")
	ErrBadRequest     = errors.New("invalid request payload")
	ErrConflict       = errors.New("resource already exists")
	ErrLinkExpired    = errors.New("link has expired")
	ErrLinkInactive   = errors.New("link is deactivated")
	ErrInternalServer = errors.New("internal server error")
)

type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(statusCode int, code, message string, details any) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

func NotFound(message string) *AppError {
	return New(http.StatusNotFound, "NOT_FOUND", message, nil)
}

func BadRequest(message string, details any) *AppError {
	return New(http.StatusBadRequest, "BAD_REQUEST", message, details)
}

func Unauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func Forbidden(message string) *AppError {
	return New(http.StatusForbidden, "FORBIDDEN", message, nil)
}

func Conflict(message string) *AppError {
	return New(http.StatusConflict, "CONFLICT", message, nil)
}

func Internal(message string) *AppError {
	return New(http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
}
