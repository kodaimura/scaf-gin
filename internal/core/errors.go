package core

import (
	"errors"
)

const (
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeConflict     = "CONFLICT"
	ErrCodeUnexpected   = "UNEXPECTED"
)

var (
	ErrBadRequest   = NewAppError("Bad request", ErrCodeBadRequest)
	ErrForbidden    = NewAppError("Forbidden", ErrCodeForbidden)
	ErrUnauthorized = NewAppError("Unauthorized access", ErrCodeUnauthorized)
	ErrNotFound     = NewAppError("Resource not found", ErrCodeNotFound)
	ErrConflict     = NewAppError("Conflict occurred", ErrCodeConflict)
	ErrUnexpected   = NewAppError("Unexpected error occurred", ErrCodeUnexpected)
)

type AppError struct {
	Message      string
	ErrorCode    string
	ErrorDetails []map[string]any
}

// Implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Get the error code of the AppError
func (e *AppError) Code() string {
	return e.ErrorCode
}

// Get additional error details
func (e *AppError) Details() []map[string]any {
	return e.ErrorDetails
}

// Check if two errors are the same based on their error code
func (e *AppError) Is(target error) bool {
	if appErr, ok := target.(*AppError); ok {
		return e.Code() == appErr.Code()
	}
	return false
}

// Create a new AppError instance
func NewAppError(message, code string) *AppError {
	return &AppError{
		Message:      message,
		ErrorCode:    code,
		ErrorDetails: []map[string]any{},
	}
}

// Create a new validation error with additional details
func NewValidationError(details []map[string]any) *AppError {
	return &AppError{
		Message:      "Bad request",
		ErrorCode:    ErrCodeBadRequest,
		ErrorDetails: details,
	}
}

// Check if the error is an AppError type
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}
