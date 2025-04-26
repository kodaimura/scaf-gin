package core

type AuthI interface {
	CreateAccessToken(payload AuthPayload) (string, error)
	CreateRefreshToken(payload AuthPayload) (string, error)
	VerifyAccessToken(token string) (AuthPayload, error)
	VerifyRefreshToken(token string) (AuthPayload, error)
	RevokeRefreshToken(token string) error
}

type AuthPayload struct {
	AccountId   int
	AccountName string
}

var Auth AuthI = &noopAuth{}

func SetAuth(a AuthI) {
	Auth = a
}

type noopAuth struct{}

func (n *noopAuth) CreateAccessToken(payload AuthPayload) (string, error)  { return "", nil }
func (n *noopAuth) CreateRefreshToken(payload AuthPayload) (string, error) { return "", nil }
func (n *noopAuth) VerifyAccessToken(token string) (AuthPayload, error)    { return AuthPayload{}, nil }
func (n *noopAuth) VerifyRefreshToken(token string) (AuthPayload, error)   { return AuthPayload{}, nil }
func (n *noopAuth) RevokeRefreshToken(token string) error                  { return nil }
