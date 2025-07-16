package core

import (
	"errors"
	"net/http"
)

const (
	ErrCodeBadRequest    = "BAD_REQUEST"
	ErrCodeUnauthorized  = "UNAUTHORIZED"
	ErrCodeForbidden     = "FORBIDDEN"
	ErrCodeNotFound      = "NOT_FOUND"
	ErrCodeConflict      = "CONFLICT"
	ErrCodeUnprocessable = "UNPROCESSABLE_ENTITY"
	ErrCodeUnexpected    = "UNEXPECTED"
)

const (
	ErrMessageBadRequest    = "Bad request"
	ErrMessageUnauthorized  = "Unauthorized access"
	ErrMessageForbidden     = "Forbidden"
	ErrMessageNotFound      = "Resource not found"
	ErrMessageConflict      = "Conflict occurred"
	ErrMessageUnprocessable = "Unprocessable entity"
	ErrMessageUnexpected    = "Unexpected error occurred"
)

var (
	ErrBadRequest    = NewAppError(ErrMessageBadRequest, ErrCodeBadRequest, http.StatusBadRequest)
	ErrUnauthorized  = NewAppError(ErrMessageUnauthorized, ErrCodeUnauthorized, http.StatusUnauthorized)
	ErrForbidden     = NewAppError(ErrMessageForbidden, ErrCodeForbidden, http.StatusForbidden)
	ErrNotFound      = NewAppError(ErrMessageNotFound, ErrCodeNotFound, http.StatusNotFound)
	ErrConflict      = NewAppError(ErrMessageConflict, ErrCodeConflict, http.StatusConflict)
	ErrUnprocessable = NewAppError(ErrMessageUnprocessable, ErrCodeUnprocessable, http.StatusUnprocessableEntity)
	ErrUnexpected    = NewAppError(ErrMessageUnexpected, ErrCodeUnexpected, http.StatusInternalServerError)
)

// AppError defines a reusable application-level error
type AppError struct {
	Message      string           `json:"message"`
	ErrorCode    string           `json:"code"`
	HTTPStatus   int              `json:"-"`
	ErrorDetails []map[string]any `json:"details,omitempty"`
	Err          error            `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap allows use of errors.Unwrap(), errors.Is(), and errors.As()
func (e *AppError) Unwrap() error {
	return e.Err
}

// Constructor: creates a new AppError without a wrapped internal error
func NewAppError(message, code string, status int) *AppError {
	return &AppError{
		Message:    message,
		ErrorCode:  code,
		HTTPStatus: status,
	}
}

// Constructor: creates a new AppError with a wrapped internal error
func WrapAppError(err error, message, code string, status int) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return &AppError{
		Message:    message,
		ErrorCode:  code,
		HTTPStatus: status,
		Err:        err,
	}
}

func NewUnexpectedError(err error) *AppError {
	return WrapAppError(err, ErrMessageUnexpected, ErrCodeUnexpected, http.StatusInternalServerError)
}

// Constructor: creates a validation error with additional detail information
func NewValidationError(details []map[string]any) *AppError {
	return &AppError{
		Message:      ErrMessageUnprocessable,
		ErrorCode:    ErrCodeUnprocessable,
		HTTPStatus:   http.StatusUnprocessableEntity,
		ErrorDetails: details,
	}
}
