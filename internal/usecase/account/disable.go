package account

import (
	"errors"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type DisableDto struct {
	Id int
}

func (uc *usecase) Disable(in DisableDto) (model.Account, error) {
	acct, err := uc.accountModule.GetByID(in.Id)
	if errors.Is(err, core.ErrNotFound) {
		return model.Account{}, core.ErrAccountNotFound
	}
	if err != nil {
		return model.Account{}, err
	}
	return uc.accountModule.Disable(acct)
}
