package auth

import (
	"gorm.io/gorm"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type SignupDto struct {
	LoginID   *string
	Email     *string
	Password  string
	FirstName string
	LastName  string
}

func (uc *usecase) Signup(in SignupDto) (model.Account, error) {
	var account model.Account
	err := uc.db.Transaction(func(tx *gorm.DB) error {
		accountModuleTx := uc.accountModule.WithTx(tx)

		loginID, err := core.ResolveLoginID(in.LoginID, in.Email)
		if err != nil {
			return err
		}
		if err := ensureUniqueLoginID(accountModuleTx, loginID, 0); err != nil {
			return err
		}
		if err := ensureUniqueEmail(accountModuleTx, in.Email, 0); err != nil {
			return err
		}

		hashed, err := hashPassword(in.Password)
		if err != nil {
			return err
		}

		account, err = accountModuleTx.Create(model.Account{
			LoginID:      loginID,
			Email:        in.Email,
			PasswordHash: hashed,
			TokenVersion: 1,
			FirstName:    in.FirstName,
			LastName:     in.LastName,
		})
		return err
	})

	return account, err
}
