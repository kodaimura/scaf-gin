package account_profile

import (
	profileModule "scaf-gin/internal/module/account_profile"
	"time"
)

// -----------------------------
// DTO（Response）
// -----------------------------

type AccountProfileResponse struct {
	AccountId   int       `json:"account_id"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToAccountProfileResponse(m profileModule.AccountProfile) AccountProfileResponse {
	return AccountProfileResponse{
		AccountId:   m.AccountId,
		DisplayName: m.DisplayName,
		Bio:         m.Bio,
		AvatarURL:   m.AvatarURL,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func ToAccountProfileResponseList(models []profileModule.AccountProfile) []AccountProfileResponse {
	res := make([]AccountProfileResponse, 0, len(models))
	for _, m := range models {
		res = append(res, ToAccountProfileResponse(m))
	}
	return res
}

// -----------------------------
// DTO（Request）
// -----------------------------

type AccountUri struct {
	AccountId int `uri:"account_id" binding:"required"`
}

type PutMeRequest struct {
	DisplayName string `json:"display_name" binding:"required"`
	Bio         string `json:"bio" binding:"omitempty"`
	AvatarURL   string `json:"avatar_url" binding:"omitempty,url"`
}
