package auth

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

type SignupResponse struct {
	Account AccountResponse `json:"account"`
}

type LoginResponse struct {
	Account     AccountResponse `json:"account"`
	AccessToken string          `json:"access_token"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

// -----------------------------
// DTO（Request）
// -----------------------------

type SignupRequest struct {
	LoginID   *string `json:"login_id" binding:"omitempty,max=255"`
	Email     *string `json:"email" binding:"omitempty,email"`
	FirstName string  `json:"first_name" binding:"required,max=100"`
	LastName  string  `json:"last_name" binding:"required,max=100"`
	Password  string  `json:"password" binding:"required,min=8,max=255"`
}

type LoginRequest struct {
	LoginID    string `json:"login_id" binding:"required,max=255"`
	Password   string `json:"password" binding:"required,max=255"`
	RememberMe bool   `json:"remember_me"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required,max=500"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=255"`
}

type PutMePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=255"`
}
