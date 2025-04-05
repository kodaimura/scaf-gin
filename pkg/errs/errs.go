package errs

import (
	"fmt"
)


/////////////////////////////////////////////////////////////////////////
type BadRequestError struct {
	Message string
}

func NewBadRequestError(message string) error {
	return BadRequestError{Message: message}
}

func (e BadRequestError) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}

/////////////////////////////////////////////////////////////////////////
type UnauthorizedError struct {}

func NewUnauthorizedError() error {
	return UnauthorizedError{}
}

func (e UnauthorizedError) Error() string {
	return "error: Authentication failed."
}

/////////////////////////////////////////////////////////////////////////
type ForbiddenError struct {}

func NewForbiddenError() error {
	return ForbiddenError{}
}

func (e ForbiddenError) Error() string {
	return "error: Permission denied to access this resource."
}

/////////////////////////////////////////////////////////////////////////
type NotFoundError struct {}

func NewNotFoundError() error {
	return NotFoundError{}
}

func (e NotFoundError) Error() string {
	return "error: Not found"
}

/////////////////////////////////////////////////////////////////////////
type ConflictError struct {}

func NewConflictError() error {
	return ConflictError{}
}

func (e ConflictError) Error() string {
	return "error: Conflict occurred due to resource uniqueness."
}

/////////////////////////////////////////////////////////////////////////
type UnexpectedError struct {
	Message string
}

func NewUnexpectedError(message string) error {
	return UnexpectedError{Message: message}
}

func (e UnexpectedError) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}