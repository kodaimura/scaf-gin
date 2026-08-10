package account

import (
	"errors"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

type UpdateDto struct {
	Id        int
	LoginID   *string
	Email     *string
	Password  *string
	FirstName string
	LastName  string
}

func (uc *usecase) Update(in UpdateDto) (accountModule.Account, error) {
	acct, err := uc.accountModule.GetByID(in.Id)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return accountModule.Account{}, core.ErrAccountNotFound
		}
		return accountModule.Account{}, err
	}

	loginID, err := core.ResolveLoginID(in.LoginID, in.Email)
	if err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueLoginID(loginID, acct.Id); err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueEmail(in.Email, acct.Id); err != nil {
		return accountModule.Account{}, err
	}

	acct.LoginID = loginID
	acct.Email = in.Email
	acct.FirstName = in.FirstName
	acct.LastName = in.LastName
	if in.Password != nil {
		hashed, err := hashPassword(*in.Password)
		if err != nil {
			return accountModule.Account{}, err
		}
		acct.PasswordHash = hashed
		acct.TokenVersion++
	}

	return uc.accountModule.Update(acct)
}
