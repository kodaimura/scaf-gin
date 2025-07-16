package helper

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"scaf-gin/internal/core"
)

func HandleGormError(err error) error {
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
	}
	if strings.Contains(err.Error(), "SQLSTATE 23505") {
		core.Logger.Debug(err.Error())
		return core.ErrConflict
	}

	core.Logger.Error(err.Error())
	return core.NewUnexpectedError(err)
}
