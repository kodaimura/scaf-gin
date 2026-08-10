package auth

import (
	"errors"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

type LoginDto struct {
	LoginID    string
	Password   string
	RememberMe bool
}

func (uc *usecase) Login(in LoginDto) (accountModule.Account, string, string, error) {
	account, err := uc.accountModule.GetByLoginID(in.LoginID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return accountModule.Account{}, "", "", core.ErrInvalidCredentials
		}
		return accountModule.Account{}, "", "", err
	}

	if !verifyPassword(account.PasswordHash, in.Password) {
		return accountModule.Account{}, "", "", core.ErrInvalidCredentials
	}
	if account.DisabledAt != nil {
		return accountModule.Account{}, "", "", core.ErrAccountDisabled
	}

	payload := core.AuthPayload{
		AccountId:    account.Id,
		LoginID:      account.LoginID,
		TokenVersion: account.TokenVersion,
	}

	accessToken, err := core.Auth.CreateAccessToken(payload)
	if err != nil {
		return accountModule.Account{}, "", "", err
	}

	refreshToken, err := core.Auth.CreateRefreshToken(payload, in.RememberMe)
	if err != nil {
		return accountModule.Account{}, "", "", err
	}
	return account, accessToken, refreshToken, nil
}
