package account

import (
	"errors"

	"scaf-gin/internal/core"
)

type DeleteDto struct {
	Id int
}

func (uc *usecase) Delete(in DeleteDto) error {
	acct, err := uc.accountModule.GetByID(in.Id)
	if errors.Is(err, core.ErrNotFound) {
		return core.ErrAccountNotFound
	}
	if err != nil {
		return err
	}
	return uc.accountModule.Delete(acct)
}
