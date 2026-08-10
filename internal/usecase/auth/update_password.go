package auth

import "scaf-gin/internal/core"

type UpdatePasswordDto struct {
	Id          int
	OldPassword string
	NewPassword string
}

func (uc *usecase) UpdatePassword(in UpdatePasswordDto) error {
	account, err := uc.accountModule.GetByID(in.Id)
	if err != nil {
		return err
	}
	if !verifyPassword(account.PasswordHash, in.OldPassword) {
		return core.ErrCurrentPasswordIncorrect
	}

	hashed, err := hashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	account.PasswordHash = hashed
	account.TokenVersion++
	_, err = uc.accountModule.Update(account)
	return err
}
