package core

import (
	"strings"

	"scaf-gin/config"
)

func ResolveLoginID(loginID *string, email *string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(config.AuthLoginIDMode)) {
	case "email":
		if email == nil || strings.TrimSpace(*email) == "" {
			return "", ErrEmailRequired
		}
		return strings.TrimSpace(*email), nil
	case "login_id":
		if loginID == nil || strings.TrimSpace(*loginID) == "" {
			return "", ErrLoginIDRequired
		}
		return strings.TrimSpace(*loginID), nil
	default:
		return "", ErrInvalidState
	}
}
