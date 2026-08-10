package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type ForgotPasswordDto struct {
	Email string
}

func (uc *usecase) ForgotPassword(in ForgotPasswordDto) error {
	account, err := uc.accountModule.GetByEmail(in.Email)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if account.DisabledAt != nil {
		return nil
	}

	now := time.Now()
	latest, err := uc.passwordResetTokenModule.FindLatestByAccountID(account.Id)
	if err != nil && !errors.Is(err, core.ErrNotFound) {
		return err
	}
	if err == nil {
		resendAfter := now.Add(-time.Minute * time.Duration(config.Current.PasswordResetResendIntervalMinutes))
		if latest.CreatedAt.After(resendAfter) {
			return nil
		}
	}

	rawToken, err := core.GenerateToken(48)
	if err != nil {
		return err
	}
	tokenHash := core.HashToken(rawToken)
	expiresAt := now.Add(time.Minute * time.Duration(config.Current.PasswordResetTokenExpiresMinutes))

	err = uc.db.Transaction(func(tx *gorm.DB) error {
		tokenModuleTx := uc.passwordResetTokenModule.WithTx(tx)
		if err := tokenModuleTx.InvalidateActiveTokens(account.Id, now); err != nil {
			return err
		}
		_, err := tokenModuleTx.Create(model.PasswordResetToken{
			AccountId: account.Id,
			TokenHash: tokenHash,
			ExpiresAt: expiresAt,
		})
		return err
	})
	if err != nil {
		return err
	}

	resetURL := buildResetURL(rawToken)
	body := buildPasswordResetMailBody(account.LastName+" "+account.FirstName, resetURL, config.Current.PasswordResetTokenExpiresMinutes)
	return uc.mailer.SendText([]string{in.Email}, "Password reset", body)
}
