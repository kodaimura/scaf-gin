package account

import (
	"errors"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type EnableDto struct {
	ID int64
}

func (uc *usecase) Enable(in EnableDto) (model.Account, error) {
	acct, err := uc.accountModule.GetByID(in.ID)
	if errors.Is(err, core.ErrNotFound) {
		return model.Account{}, core.ErrAccountNotFound
	}
	if err != nil {
		return model.Account{}, err
	}
	return uc.accountModule.Enable(acct)
}
