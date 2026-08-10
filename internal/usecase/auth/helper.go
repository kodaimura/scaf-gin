package auth

import (
	"errors"
	"net/url"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"scaf-gin/config"
	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

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
