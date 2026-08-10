package core

type AuthI interface {
	CreateAccessToken(payload AuthPayload) (string, error)
	CreateRefreshToken(payload AuthPayload, rememberMe bool) (string, error)
	VerifyAccessToken(token string) (AuthPayload, error)
	VerifyRefreshToken(token string) (AuthPayload, error)
	RevokeRefreshToken(token string) error
}

type AuthPayload struct {
	AccountId    int    `json:"account_id"`
	LoginID      string `json:"login_id"`
	TokenVersion int    `json:"token_version"`
}
