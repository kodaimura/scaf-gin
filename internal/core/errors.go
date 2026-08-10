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

	ErrCodeAccountDisabled          = "ACCOUNT_DISABLED"
	ErrCodeAccountNotFound          = "ACCOUNT_NOT_FOUND"
	ErrCodeAuthInvalid              = "AUTH_INVALID"
	ErrCodeAuthInvalidPayload       = "AUTH_INVALID_PAYLOAD"
	ErrCodeAuthInvalidSubject       = "AUTH_INVALID_SUBJECT"
	ErrCodeAuthInvalidType          = "AUTH_INVALID_TYPE"
	ErrCodeAuthMissing              = "AUTH_MISSING"
	ErrCodeAuthNotFound             = "AUTH_NOT_FOUND"
	ErrCodeAuthTokenRevoked         = "AUTH_TOKEN_REVOKED"
	ErrCodeCurrentPasswordIncorrect = "CURRENT_PASSWORD_INCORRECT"
	ErrCodeInvalidCredentials       = "INVALID_CREDENTIALS"
	ErrCodeMalformedToken           = "MALFORMED_TOKEN"
	ErrCodeRefreshInvalid           = "REFRESH_INVALID"
	ErrCodeRefreshInvalidPayload    = "REFRESH_INVALID_PAYLOAD"
	ErrCodeRefreshInvalidType       = "REFRESH_INVALID_TYPE"
	ErrCodeRefreshMissing           = "REFRESH_MISSING"
	ErrCodeTokenAlreadyUsed         = "TOKEN_ALREADY_USED"
	ErrCodeTokenExpired             = "TOKEN_EXPIRED"
	ErrCodeTokenInvalid             = "TOKEN_INVALID"
	ErrCodeInvalidState             = "INVALID_STATE"
	ErrCodeLoginIDAlreadyExists     = "LOGIN_ID_ALREADY_EXISTS"
	ErrCodeEmailAlreadyExists       = "EMAIL_ALREADY_EXISTS"
	ErrCodeEmailRequired            = "EMAIL_REQUIRED"
	ErrCodeLoginIDRequired          = "LOGIN_ID_REQUIRED"
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

	ErrAccountDisabled          = NewAppError("Account is disabled", ErrCodeAccountDisabled, http.StatusUnauthorized)
	ErrAccountNotFound          = NewAppError("Account not found", ErrCodeAccountNotFound, http.StatusNotFound)
	ErrAuthInvalid              = NewAppError("Invalid access token", ErrCodeAuthInvalid, http.StatusUnauthorized)
	ErrAuthInvalidPayload       = NewAppError("Invalid access token payload", ErrCodeAuthInvalidPayload, http.StatusUnauthorized)
	ErrAuthInvalidSubject       = NewAppError("Invalid access token subject", ErrCodeAuthInvalidSubject, http.StatusUnauthorized)
	ErrAuthInvalidType          = NewAppError("Invalid access token type", ErrCodeAuthInvalidType, http.StatusUnauthorized)
	ErrAuthMissing              = NewAppError("Missing access token", ErrCodeAuthMissing, http.StatusUnauthorized)
	ErrAuthNotFound             = NewAppError("Authenticated account not found", ErrCodeAuthNotFound, http.StatusUnauthorized)
	ErrAuthTokenRevoked         = NewAppError("Authentication token was revoked", ErrCodeAuthTokenRevoked, http.StatusUnauthorized)
	ErrCurrentPasswordIncorrect = NewAppError("Current password is incorrect", ErrCodeCurrentPasswordIncorrect, http.StatusUnauthorized)
	ErrInvalidCredentials       = NewAppError("Invalid credentials", ErrCodeInvalidCredentials, http.StatusUnauthorized)
	ErrMalformedToken           = NewAppError("Malformed token", ErrCodeMalformedToken, http.StatusUnauthorized)
	ErrRefreshInvalid           = NewAppError("Invalid refresh token", ErrCodeRefreshInvalid, http.StatusUnauthorized)
	ErrRefreshInvalidPayload    = NewAppError("Invalid refresh token payload", ErrCodeRefreshInvalidPayload, http.StatusUnauthorized)
	ErrRefreshInvalidType       = NewAppError("Invalid refresh token type", ErrCodeRefreshInvalidType, http.StatusUnauthorized)
	ErrRefreshMissing           = NewAppError("Missing refresh token", ErrCodeRefreshMissing, http.StatusUnauthorized)
	ErrTokenAlreadyUsed         = NewAppError("Token already used", ErrCodeTokenAlreadyUsed, http.StatusBadRequest)
	ErrTokenExpired             = NewAppError("Token expired", ErrCodeTokenExpired, http.StatusBadRequest)
	ErrTokenInvalid             = NewAppError("Invalid token", ErrCodeTokenInvalid, http.StatusBadRequest)
	ErrInvalidState             = NewAppError("Invalid state", ErrCodeInvalidState, http.StatusBadRequest)
	ErrLoginIDAlreadyExists     = NewAppError("Login ID already exists", ErrCodeLoginIDAlreadyExists, http.StatusConflict)
	ErrEmailAlreadyExists       = NewAppError("Email already exists", ErrCodeEmailAlreadyExists, http.StatusConflict)
	ErrEmailRequired            = NewAppError("Email is required", ErrCodeEmailRequired, http.StatusBadRequest)
	ErrLoginIDRequired          = NewAppError("Login ID is required", ErrCodeLoginIDRequired, http.StatusBadRequest)
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
