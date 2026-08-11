package account

import (
	"scaf-gin/internal/model"
	helper "scaf-gin/internal/usecase"
)

type CreateDto struct {
	LoginID   *string
	Email     *string
	Password  string
	FirstName string
	LastName  string
}

func (uc *usecase) Create(in CreateDto) (model.Account, error) {
	loginID, err := helper.ResolveLoginID(in.LoginID, in.Email)
	if err != nil {
		return model.Account{}, err
	}
	if err := uc.ensureUniqueLoginID(loginID, 0); err != nil {
		return model.Account{}, err
	}
	if err := uc.ensureUniqueEmail(in.Email, 0); err != nil {
		return model.Account{}, err
	}

	hashed, err := hashPassword(in.Password)
	if err != nil {
		return model.Account{}, err
	}

	return uc.accountModule.Create(model.Account{
		LoginID:      loginID,
		Email:        in.Email,
		PasswordHash: hashed,
		TokenVersion: 1,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
	})
}
