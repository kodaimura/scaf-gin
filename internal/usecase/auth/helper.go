package auth

import (
	"errors"
	"net/url"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	"scaf-gin/internal/module"
)

func ensureUniqueLoginID(accountModule module.AccountModule, loginID string, exceptID int64) error {
	existing, err := accountModule.GetByLoginID(loginID)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.ID != exceptID {
		return core.ErrLoginIDAlreadyExists
	}
	return nil
}

func ensureUniqueEmail(accountModule module.AccountModule, email *string, exceptID int64) error {
	if email == nil || *email == "" {
		return nil
	}
	existing, err := accountModule.GetByEmail(*email)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.ID != exceptID {
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

func buildResetURL(token string) string {
	resetURL, err := url.Parse(config.Current.PasswordResetURLBase)
	if err != nil {
		return config.Current.PasswordResetURLBase + "?token=" + url.QueryEscape(token)
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
