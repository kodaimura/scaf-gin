package service

import (
	"goscaf/internal/core"
)

func handleError(err error) error {
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