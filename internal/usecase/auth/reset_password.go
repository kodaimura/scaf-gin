package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type ResetPasswordDto struct {
	Token       string
	NewPassword string
}

func (uc *usecase) ResetPassword(in ResetPasswordDto) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		accountModuleTx := uc.accountModule.WithTx(tx)
		tokenModuleTx := uc.passwordResetTokenModule.WithTx(tx)

		token, err := tokenModuleTx.GetByHashForUpdate(core.HashToken(in.Token))
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				return core.ErrTokenInvalid
			}
			return err
		}
		if err := validateResetToken(token); err != nil {
			return err
		}

		account, err := accountModuleTx.GetByID(token.AccountID)
		if err != nil {
			if errors.Is(err, core.ErrNotFound) {
				return core.ErrAccountNotFound
			}
			return err
		}

		hashed, err := hashPassword(in.NewPassword)
		if err != nil {
			return err
		}

		now := time.Now()
		account.PasswordHash = hashed
		account.TokenVersion++
		token.UsedAt = &now

		if _, err := accountModuleTx.Update(account); err != nil {
			return err
		}
		_, err = tokenModuleTx.Update(token)
		return err
	})
}

func validateResetToken(token model.PasswordResetToken) error {
	if token.UsedAt != nil {
		return core.ErrTokenAlreadyUsed
	}
	if !token.ExpiresAt.After(time.Now()) {
		return core.ErrTokenExpired
	}
	return nil
}
