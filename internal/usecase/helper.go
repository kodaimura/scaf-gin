package usecase

import (
	"strings"

	"scaf-gin/config"
	"scaf-gin/internal/core"
)

func ResolveLoginID(loginID *string, email *string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(config.Current.AuthLoginIDMode)) {
	case "email":
		if email == nil || strings.TrimSpace(*email) == "" {
			return "", core.ErrEmailRequired
		}
		return strings.TrimSpace(*email), nil
	case "login_id":
		if loginID == nil || strings.TrimSpace(*loginID) == "" {
			return "", core.ErrLoginIDRequired
		}
		return strings.TrimSpace(*loginID), nil
	default:
		return "", core.ErrInvalidState
	}
}
