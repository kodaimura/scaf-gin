package auth

import (
	"errors"

	"scaf-gin/internal/core"
)

func (uc *usecase) AuthorizeAccessToken(payload core.AuthPayload) (int64, error) {
	if payload.AccountID == 0 || payload.TokenVersion == 0 {
		return 0, core.ErrAuthInvalidPayload
	}

	account, err := uc.accountModule.GetByID(payload.AccountID)
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

	return account.ID, nil
}
