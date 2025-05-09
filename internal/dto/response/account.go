package response

import (
	"scaf-gin/internal/model"
	"time"
)

// ============================
// Account
// ============================

type Account struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FromModelAccount(m model.Account) Account {
	return Account{
		Id:        m.Id,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func FromModelAccountList(models []model.Account) []Account {
	res := make([]Account, 0, len(models))
	for _, m := range models {
		res = append(res, FromModelAccount(m))
	}
	return res
}

// ============================
// Login
// ============================

type Login struct {
	AccessToken      string  `json:"access_token"`
	RefreshToken     string  `json:"refresh_token"`
	AccessExpiresIn  int     `json:"access_expires_in"`
	RefreshExpiresIn int     `json:"refresh_expires_in"`
	Account          Account `json:"account"`
}

// ============================
// Refresh
// ============================

type Refresh struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
