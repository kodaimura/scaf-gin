package account

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"scaf-gin/internal/core"
)

func (uc *usecase) ensureUniqueLoginID(loginID string, exceptID int64) error {
	existing, err := uc.accountModule.GetByLoginID(loginID)
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

func (uc *usecase) ensureUniqueEmail(email *string, exceptID int64) error {
	if email == nil || *email == "" {
		return nil
	}
	existing, err := uc.accountModule.GetByEmail(*email)
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
