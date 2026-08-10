package auth

import (
	"errors"

	"scaf-gin/internal/core"
)

func (uc *usecase) Refresh(refreshToken string) (core.AuthPayload, string, error) {
	if refreshToken == "" {
		return core.AuthPayload{}, "", core.ErrRefreshMissing
	}
	payload, err := uc.auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		return core.AuthPayload{}, "", core.ErrRefreshInvalid
	}

	account, err := uc.accountModule.GetByID(payload.AccountId)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.AuthPayload{}, "", core.ErrAuthNotFound
		}
		return core.AuthPayload{}, "", err
	}
	if account.DisabledAt != nil {
		return core.AuthPayload{}, "", core.ErrAccountDisabled
	}
	if payload.TokenVersion != account.TokenVersion {
		return core.AuthPayload{}, "", core.ErrAuthTokenRevoked
	}

	accessToken, err := uc.auth.CreateAccessToken(core.AuthPayload{
		AccountId:    account.Id,
		TokenVersion: account.TokenVersion,
	})

	return payload, accessToken, err
}
