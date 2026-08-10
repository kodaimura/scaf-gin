package module

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"scaf-gin/internal/core"
)

func handleGormError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return core.ErrNotFound
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return core.ErrConflict
	}
	if strings.Contains(err.Error(), "SQLSTATE 23505") {
		return core.ErrConflict
	}

	return core.NewUnexpectedError(err)
}
