package account

import (
	"errors"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

type EnableDto struct {
	Id int
}

func (uc *usecase) Enable(in EnableDto) (accountModule.Account, error) {
	acct, err := uc.accountModule.GetByID(in.Id)
	if errors.Is(err, core.ErrNotFound) {
		return accountModule.Account{}, core.ErrAccountNotFound
	}
	if err != nil {
		return accountModule.Account{}, err
	}
	return uc.accountModule.Enable(acct)
}
