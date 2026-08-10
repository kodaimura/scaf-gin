package service

import (
	"errors"

	"gorm.io/gorm"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

func ValidateAccessTokenAccount(db *gorm.DB, payload core.AuthPayload) (int, error) {
	if payload.AccountId == 0 || payload.TokenVersion == 0 {
		return 0, core.ErrAuthInvalidPayload
	}

	account, err := accountModule.NewModule(db).GetByID(payload.AccountId)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return 0, core.ErrAuthNotFound
		}
		return 0, err
	}

	if account.DisabledAt != nil {
		return 0, core.ErrAccountDisabled
	}

	if payload.TokenVersion != account.TokenVersion {
		return 0, core.ErrAuthTokenRevoked
	}

	return account.Id, nil
}
