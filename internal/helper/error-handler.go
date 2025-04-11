package helper

import (
	"goscaf/internal/core"
)

// HandleError handles the given error, checks if it is a custom application error,
// and returns the appropriate error. If the error is not a recognized app error, 
// it logs the error and returns a generic unexpected error.
func HandleError(err error) error {
	if err == nil {
		return nil
	}
	if core.IsAppError(err) {
		return err
	} else {
		core.Logger.Error(err.Error())
		return core.ErrUnexpected
	}
}