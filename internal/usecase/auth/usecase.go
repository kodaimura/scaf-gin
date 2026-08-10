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
	profileModule "scaf-gin/internal/module/account_profile"
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
	db                        *gorm.DB
	accountService            accountModule.Service
	accountProfileService     profileModule.Service
	passwordResetTokenService passwordResetTokenModule.Service
}

func NewUsecase(
	db *gorm.DB,
	accountService accountModule.Service,
	accountProfileService profileModule.Service,
	passwordResetTokenService passwordResetTokenModule.Service,
) Usecase {
	return &usecase{
		db:                        db,
		accountService:            accountService,
		accountProfileService:     accountProfileService,
		passwordResetTokenService: passwordResetTokenService,
	}
}

func (uc *usecase) Signup(in SignupDto) (accountModule.Account, error) {
	var account accountModule.Account
	err := uc.db.Transaction(func(tx *gorm.DB) error {
		loginID, err := core.ResolveLoginID(in.LoginID, in.Email)
		if err != nil {
			return err
		}
		if err := uc.ensureUniqueLoginID(loginID, 0, tx); err != nil {
			return err
		}
		if err := uc.ensureUniqueEmail(in.Email, 0, tx); err != nil {
			return err
		}

		hashed, err := hashPassword(in.Password)
		if err != nil {
			return err
		}

		account, err = uc.accountService.CreateOne(accountModule.Account{
			LoginID:      loginID,
			Email:        in.Email,
			PasswordHash: string(hashed),
			TokenVersion: 1,
			FirstName:    in.FirstName,
			LastName:     in.LastName,
		}, tx)

		if err != nil {
			return err
		}

		_, err = uc.accountProfileService.CreateOne(profileModule.AccountProfile{
			AccountId:   account.Id,
			DisplayName: account.LastName + " " + account.FirstName,
			Bio:         "",
			AvatarURL:   "",
		}, tx)
		return err
	})

	return account, err
}

func (uc *usecase) Login(in LoginDto) (accountModule.Account, string, string, error) {
	acct, err := uc.accountService.GetByLoginID(in.LoginID, uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return accountModule.Account{}, "", "", core.ErrInvalidCredentials
		}
		return accountModule.Account{}, "", "", err
	}

	if !verifyPassword(acct.PasswordHash, in.Password) {
		return accountModule.Account{}, "", "", core.ErrInvalidCredentials
	}
	if acct.DisabledAt != nil {
		return accountModule.Account{}, "", "", core.ErrAccountDisabled
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:    acct.Id,
		LoginID:      acct.LoginID,
		TokenVersion: acct.TokenVersion,
	})
	if err != nil {
		return accountModule.Account{}, "", "", err
	}

	refreshToken, err := core.Auth.CreateRefreshToken(core.AuthPayload{
		AccountId:    acct.Id,
		LoginID:      acct.LoginID,
		TokenVersion: acct.TokenVersion,
	}, in.RememberMe)
	if err != nil {
		return accountModule.Account{}, "", "", err
	}
	return acct, accessToken, refreshToken, nil
}

func (uc *usecase) Refresh(refreshToken string) (core.AuthPayload, string, error) {
	if refreshToken == "" {
		return core.AuthPayload{}, "", core.ErrRefreshMissing
	}
	payload, err := core.Auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		return core.AuthPayload{}, "", core.ErrRefreshInvalid
	}

	acct, err := uc.accountService.GetOne(accountModule.Account{Id: payload.AccountId}, uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.AuthPayload{}, "", core.ErrAuthNotFound
		}
		return core.AuthPayload{}, "", err
	}
	if acct.DisabledAt != nil {
		return core.AuthPayload{}, "", core.ErrAccountDisabled
	}
	if payload.TokenVersion != acct.TokenVersion {
		return core.AuthPayload{}, "", core.ErrAuthTokenRevoked
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:    acct.Id,
		LoginID:      acct.LoginID,
		TokenVersion: acct.TokenVersion,
	})

	return payload, accessToken, err
}

func (uc *usecase) ForgotPassword(in ForgotPasswordDto) error {
	acct, err := uc.accountService.GetByEmail(in.Email, uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if acct.DisabledAt != nil {
		return nil
	}

	now := time.Now()
	latest, err := uc.passwordResetTokenService.FindLatestByAccountId(acct.Id, uc.db)
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
		if err := uc.passwordResetTokenService.InvalidateActiveTokens(acct.Id, now, tx); err != nil {
			return err
		}
		_, err := uc.passwordResetTokenService.CreateOne(passwordResetTokenModule.PasswordResetToken{
			AccountId: acct.Id,
			TokenHash: tokenHash,
			ExpiresAt: expiresAt,
		}, tx)
		return err
	})
	if err != nil {
		return err
	}

	resetURL := buildResetURL(rawToken)
	body := buildPasswordResetMailBody(acct.LastName+" "+acct.FirstName, resetURL, config.PasswordResetTokenExpiresMinutes)
	return core.Mailer.SendText([]string{in.Email}, "Password reset", body)
}

func (uc *usecase) VerifyResetPasswordToken(in VerifyResetPasswordTokenDto) error {
	token, err := uc.passwordResetTokenService.GetByHash(core.HashToken(in.Token), uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.ErrTokenInvalid
		}
		return err
	}
	return validateResetToken(token)
}

func (uc *usecase) ResetPassword(in ResetPasswordDto) error {
	token, err := uc.passwordResetTokenService.GetByHash(core.HashToken(in.Token), uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return core.ErrTokenInvalid
		}
		return err
	}
	if err := validateResetToken(token); err != nil {
		return err
	}

	acct, err := uc.accountService.GetOne(accountModule.Account{Id: token.AccountId}, uc.db)
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
	acct.PasswordHash = hashed
	acct.TokenVersion++
	token.UsedAt = &now

	return uc.db.Transaction(func(tx *gorm.DB) error {
		if _, err := uc.accountService.UpdateOne(acct, tx); err != nil {
			return err
		}
		_, err := uc.passwordResetTokenService.UpdateOne(token, tx)
		return err
	})
}

func (uc *usecase) UpdatePassword(in UpdatePasswordDto) error {
	acct, err := uc.accountService.GetOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
	if err != nil {
		return err
	}
	if !verifyPassword(acct.PasswordHash, in.OldPassword) {
		return core.ErrCurrentPasswordIncorrect
	}

	hashed, err := hashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	acct.PasswordHash = string(hashed)
	acct.TokenVersion++
	_, err = uc.accountService.UpdateOne(acct, uc.db)
	return err
}

func (uc *usecase) ensureUniqueLoginID(loginID string, exceptID int, db *gorm.DB) error {
	existing, err := uc.accountService.GetByLoginID(loginID, db)
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

func (uc *usecase) ensureUniqueEmail(email *string, exceptID int, db *gorm.DB) error {
	if email == nil || *email == "" {
		return nil
	}
	existing, err := uc.accountService.GetByEmail(*email, db)
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
