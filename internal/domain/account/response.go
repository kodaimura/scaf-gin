package account

import (
	"time"
)

// ============================
// Account
// ============================

type AccountResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromModelAccount(m Account) AccountResponse {
	return AccountResponse{
		Id:        m.Id,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromModelAccountList(models []Account) []AccountResponse {
	res := make([]AccountResponse, 0, len(models))
	for _, m := range models {
		res = append(res, FromModelAccount(m))
	}
	return res
}

// ============================
// Login
// ============================

type LoginResponse struct {
	AccessToken      string  `json:"access_token"`
	RefreshToken     string  `json:"refresh_token"`
	AccessExpiresIn  int     `json:"access_expires_in"`
	RefreshExpiresIn int     `json:"refresh_expires_in"`
	Account          AccountResponse `json:"account"`
}

// ============================
// Refresh
// ============================

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
