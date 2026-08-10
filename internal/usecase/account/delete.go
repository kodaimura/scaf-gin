package account

import (
	"errors"

	"scaf-gin/internal/core"
)

type DeleteDto struct {
	ID int
}

func (uc *usecase) Delete(in DeleteDto) error {
	acct, err := uc.accountModule.GetByID(in.ID)
	if errors.Is(err, core.ErrNotFound) {
		return core.ErrAccountNotFound
	}
	if err != nil {
		return err
	}
	return uc.accountModule.Delete(acct)
}
