package impl

import (
	"errors"
	"gorm.io/gorm"

	"scaf-gin/internal/core"
)

func handleGormError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		core.Logger.Debug(err.Error())
		return core.ErrNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		core.Logger.Debug(err.Error())
		return core.ErrConflict
	} else {
		core.Logger.Error(err.Error())
		return core.ErrUnexpected
	}
}