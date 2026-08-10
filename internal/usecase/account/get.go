package account

import (
	"errors"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

type GetDto struct {
	Id int
}

func (uc *usecase) Get(in GetDto) (accountModule.Account, error) {
	acct, err := uc.accountModule.GetByID(in.Id)
	if errors.Is(err, core.ErrNotFound) {
		return accountModule.Account{}, core.ErrAccountNotFound
	}
	return acct, err
}
