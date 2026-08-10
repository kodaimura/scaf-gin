package auth

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
	passwordResetTokenModule "scaf-gin/internal/module/password_reset_token"
)

type Usecase interface {
	Signup(in SignupDto) (accountModule.Account, error)
	Login(in LoginDto) (accountModule.Account, string, string, error)
	Refresh(refreshToken string) (core.AuthPayload, string, error)
	ForgotPassword(in ForgotPasswordDto) error
	VerifyResetPasswordToken(in VerifyResetPasswordTokenDto) error
	ResetPassword(in ResetPasswordDto) error
	UpdatePassword(in UpdatePasswordDto) error
}

type usecase struct {
	db                       *gorm.DB
	accountModule            accountModule.Module
	passwordResetTokenModule passwordResetTokenModule.Module
}

func NewUsecase(
	db *gorm.DB,
	accountModule accountModule.Module,
	passwordResetTokenModule passwordResetTokenModule.Module,
) Usecase {
	return &usecase{
		db:                       db,
		accountModule:            accountModule,
		passwordResetTokenModule: passwordResetTokenModule,
	}
}

func (uc *usecase) Signup(in SignupDto) (accountModule.Account, error) {
	var account accountModule.Account
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

		account, err = accountModuleTx.Create(accountModule.Account{
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

func (uc *usecase) Login(in LoginDto) (accountModule.Account, string, string, error) {
	account, err := uc.accountModule.GetByLoginID(in.LoginID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return accountModule.Account{}, "", "", core.ErrInvalidCredentials
		}
		return accountModule.Account{}, "", "", err
	}

	if !verifyPassword(account.PasswordHash, in.Password) {
		return accountModule.Account{}, "", "", core.ErrInvalidCredentials
	}
	if account.DisabledAt != nil {
		return accountModule.Account{}, "", "", core.ErrAccountDisabled
	}

	payload := core.AuthPayload{
		AccountId:    account.Id,
		LoginID:      account.LoginID,
		TokenVersion: account.TokenVersion,
	}

	accessToken, err := core.Auth.CreateAccessToken(payload)
	if err != nil {
		return accountModule.Account{}, "", "", err
	}

	refreshToken, err := core.Auth.CreateRefreshToken(payload, in.RememberMe)
	if err != nil {
		return accountModule.Account{}, "", "", err
	}
	return account, accessToken, refreshToken, nil
}

func (uc *usecase) Refresh(refreshToken string) (core.AuthPayload, string, error) {
	if refreshToken == "" {
		return core.AuthPayload{}, "", core.ErrRefreshMissing
	}
	payload, err := core.Auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		return core.AuthPayload{}, "", core.ErrRefreshInvalid
	}

	account, err := uc.accountModule.GetByID(payload.AccountId)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.AuthPayload{}, "", core.ErrAuthNotFound
		}
		return core.AuthPayload{}, "", err
	}
	if account.DisabledAt != nil {
		return core.AuthPayload{}, "", core.ErrAccountDisabled
	}
	if payload.TokenVersion != account.TokenVersion {
		return core.AuthPayload{}, "", core.ErrAuthTokenRevoked
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:    account.Id,
		LoginID:      account.LoginID,
		TokenVersion: account.TokenVersion,
	})

	return payload, accessToken, err
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
		resendAfter := now.Add(-time.Minute * time.Duration(config.PasswordResetResendIntervalMinutes))
		if latest.CreatedAt.After(resendAfter) {
			return nil
		}
	}

	rawToken, err := core.GenerateToken(48)
	if err != nil {
		return err
	}
	tokenHash := core.HashToken(rawToken)
	expiresAt := now.Add(time.Minute * time.Duration(config.PasswordResetTokenExpiresMinutes))

	err = uc.db.Transaction(func(tx *gorm.DB) error {
		tokenModuleTx := uc.passwordResetTokenModule.WithTx(tx)
		if err := tokenModuleTx.InvalidateActiveTokens(account.Id, now); err != nil {
			return err
		}
		_, err := tokenModuleTx.Create(passwordResetTokenModule.PasswordResetToken{
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
	body := buildPasswordResetMailBody(account.LastName+" "+account.FirstName, resetURL, config.PasswordResetTokenExpiresMinutes)
	return core.Mailer.SendText([]string{in.Email}, "Password reset", body)
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

func (uc *usecase) ResetPassword(in ResetPasswordDto) error {
	token, err := uc.passwordResetTokenModule.GetByHash(core.HashToken(in.Token))
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.ErrTokenInvalid
		}
		return err
	}
	if err := validateResetToken(token); err != nil {
		return err
	}

	account, err := uc.accountModule.GetByID(token.AccountId)
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

	return uc.db.Transaction(func(tx *gorm.DB) error {
		accountModuleTx := uc.accountModule.WithTx(tx)
		tokenModuleTx := uc.passwordResetTokenModule.WithTx(tx)
		if _, err := accountModuleTx.Update(account); err != nil {
			return err
		}
		_, err := tokenModuleTx.Update(token)
		return err
	})
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

func ensureUniqueLoginID(module accountModule.Module, loginID string, exceptID int) error {
	existing, err := module.GetByLoginID(loginID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.Id != exceptID {
		return core.ErrLoginIDAlreadyExists
	}
	return nil
}

func ensureUniqueEmail(module accountModule.Module, email *string, exceptID int) error {
	if email == nil || *email == "" {
		return nil
	}
	existing, err := module.GetByEmail(*email)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.Id != exceptID {
		return core.ErrEmailAlreadyExists
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(hashedPassword, plainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return false
	}
	return true
}

func validateResetToken(token passwordResetTokenModule.PasswordResetToken) error {
	if token.UsedAt != nil {
		return core.ErrTokenAlreadyUsed
	}
	if !token.ExpiresAt.After(time.Now()) {
		return core.ErrTokenExpired
	}
	return nil
}

func buildResetURL(token string) string {
	resetURL, err := url.Parse(config.PasswordResetURLBase)
	if err != nil {
		return config.PasswordResetURLBase + "?token=" + url.QueryEscape(token)
	}
	query := resetURL.Query()
	query.Set("token", token)
	resetURL.RawQuery = query.Encode()
	return resetURL.String()
}

func buildPasswordResetMailBody(name string, resetURL string, expiresMinutes int) string {
	return "Hello " + name + ",\n\n" +
		"We received a request to reset your password.\n" +
		"Open the link below to set a new password.\n\n" +
		resetURL + "\n\n" +
		"This link expires in " + strconv.Itoa(expiresMinutes) + " minutes.\n" +
		"If you did not request this, you can ignore this email."
}
