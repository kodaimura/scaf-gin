package auth

import (
	"errors"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type LoginDto struct {
	LoginID    string
	Password   string
	RememberMe bool
}

func (uc *usecase) Login(in LoginDto) (model.Account, string, string, error) {
	account, err := uc.accountModule.GetByLoginID(in.LoginID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return model.Account{}, "", "", core.ErrInvalidCredentials
		}
		return model.Account{}, "", "", err
	}

	if !verifyPassword(account.PasswordHash, in.Password) {
		return model.Account{}, "", "", core.ErrInvalidCredentials
	}
	if account.DisabledAt != nil {
		return model.Account{}, "", "", core.ErrAccountDisabled
	}

	payload := core.AuthPayload{
		AccountId:    account.Id,
		LoginID:      account.LoginID,
		TokenVersion: account.TokenVersion,
	}

	accessToken, err := core.Auth.CreateAccessToken(payload)
	if err != nil {
		return model.Account{}, "", "", err
	}

	refreshToken, err := core.Auth.CreateRefreshToken(payload, in.RememberMe)
	if err != nil {
		return model.Account{}, "", "", err
	}
	return account, accessToken, refreshToken, nil
}
