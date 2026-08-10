package account

import (
	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

type CreateDto struct {
	LoginID   *string
	Email     *string
	Password  string
	FirstName string
	LastName  string
}

func (uc *usecase) Create(in CreateDto) (accountModule.Account, error) {
	loginID, err := core.ResolveLoginID(in.LoginID, in.Email)
	if err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueLoginID(loginID, 0); err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueEmail(in.Email, 0); err != nil {
		return accountModule.Account{}, err
	}

	hashed, err := hashPassword(in.Password)
	if err != nil {
		return accountModule.Account{}, err
	}

	return uc.accountModule.Create(accountModule.Account{
		LoginID:      loginID,
		Email:        in.Email,
		PasswordHash: hashed,
		TokenVersion: 1,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
	})
}
