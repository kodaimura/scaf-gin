package account

import (
	"errors"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type GetDto struct {
	ID int64
}

func (uc *usecase) Get(in GetDto) (model.Account, error) {
	acct, err := uc.accountModule.GetByID(in.ID)
	if errors.Is(err, core.ErrNotFound) {
		return model.Account{}, core.ErrAccountNotFound
	}
	return acct, err
}
