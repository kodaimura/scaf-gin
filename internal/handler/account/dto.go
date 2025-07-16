package account

import (
	"time"

	accountModule "scaf-gin/internal/module/account"
)

// -----------------------------
// DTO（Response）
// -----------------------------

type AccountResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToAccountResponse(m accountModule.Account) AccountResponse {
	return AccountResponse{
		Id:        m.Id,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func ToAccountResponseList(models []accountModule.Account) []AccountResponse {
	res := make([]AccountResponse, 0, len(models))
	for _, m := range models {
		res = append(res, ToAccountResponse(m))
	}
	return res
}

// -----------------------------
// DTO（Request）
// -----------------------------

type PutMeRequest struct {
	Name string `json:"name" binding:"required"`
}
