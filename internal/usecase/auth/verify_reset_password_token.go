package auth

import (
	"errors"

	"scaf-gin/internal/core"
)

type VerifyResetPasswordTokenDto struct {
	Token string
}

func (uc *usecase) VerifyResetPasswordToken(in VerifyResetPasswordTokenDto) error {
	token, err := uc.passwordResetTokenModule.GetByHash(core.HashToken(in.Token))
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.ErrTokenInvalid
		}
		return err
	}
	return validateResetToken(token)
}
