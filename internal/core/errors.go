package core

import (
	"errors"
)

var (
	ErrBadRequest   = NewAppError("Bad request", "BAD_REQUEST")
	ErrForbidden    = NewAppError("Forbidden", "FORBIDDEN")
	ErrUnauthorized = NewAppError("Unauthorized access", "UNAUTHORIZED")
	ErrNotFound     = NewAppError("Resource not found", "NOT_FOUND")
	ErrConflict     = NewAppError("Conflict occurred", "CONFLICT")
	ErrUnexpected   = NewAppError("Unexpected error occurred", "UNEXPECTED")
)

type AppError struct {
	Message      string
	ErrorCode    string
	ErrorDetails []map[string]interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Code() string {
	return e.ErrorCode
}

func (e *AppError) Details() []map[string]interface{} {
	return e.ErrorDetails
}

func (e *AppError) Is(target error) bool {
	if appErr, ok := target.(*AppError); ok {
		return e.Code() == appErr.Code()
	}
	return false
}

func NewAppError(message, code string) *AppError {
	return &AppError{
		Message:      message,
		ErrorCode:    code,
		ErrorDetails: []map[string]interface{}{},
	}
}

func NewValidationError(details []map[string]interface{}) *AppError {
	return &AppError{
		Message:      "Bad request",
		ErrorCode:    "BAD_REQUEST",
		ErrorDetails: details,
	}
}

func IsAppError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return true
	}
	return false
}
