package account

import (
	"time"

	"scaf-gin/internal/model"
)

// -----------------------------
// DTO（Response）
// -----------------------------

type AccountResponse struct {
	ID         int64      `json:"id"`
	Email      *string    `json:"email"`
	LoginID    string     `json:"login_id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	DisabledAt *time.Time `json:"disabled_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func ToAccountResponse(m model.Account) AccountResponse {
	return AccountResponse{
		ID:         m.ID,
		Email:      m.Email,
		LoginID:    m.LoginID,
		FirstName:  m.FirstName,
		LastName:   m.LastName,
		DisabledAt: m.DisabledAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func ToAccountResponseList(models []model.Account) []AccountResponse {
	res := make([]AccountResponse, 0, len(models))
	for _, m := range models {
		res = append(res, ToAccountResponse(m))
	}
	return res
}

// -----------------------------
// DTO（Request）
// -----------------------------

type PostAccountRequest struct {
	LoginID   *string `json:"login_id" binding:"omitempty,max=255"`
	Email     *string `json:"email" binding:"omitempty,email"`
	FirstName string  `json:"first_name" binding:"required,max=100"`
	LastName  string  `json:"last_name" binding:"required,max=100"`
	Password  string  `json:"password" binding:"required,min=8,max=255"`
}

type PutAccountRequest struct {
	LoginID   *string `json:"login_id" binding:"omitempty,max=255"`
	Email     *string `json:"email" binding:"omitempty,email"`
	FirstName string  `json:"first_name" binding:"required,max=100"`
	LastName  string  `json:"last_name" binding:"required,max=100"`
	Password  *string `json:"password" binding:"omitempty,min=8,max=255"`
}

type GetAccountsResponse struct {
	Accounts []AccountResponse `json:"accounts"`
}

type GetAccountResponse struct {
	Account AccountResponse `json:"account"`
}

type PostAccountResponse struct {
	Account AccountResponse `json:"account"`
}

type PutAccountResponse struct {
	Account AccountResponse `json:"account"`
}

type PutAccountDisableResponse struct {
	Account AccountResponse `json:"account"`
}

type PutAccountEnableResponse struct {
	Account AccountResponse `json:"account"`
}
