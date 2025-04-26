package core

type AuthI interface {
	GenerateAccessToken(payload AuthPayload) (string, error)
	GenerateRefreshToken(payload AuthPayload) (string, error)
	ValidateAccessToken(token string) (AuthPayload, error)
	ValidateRefreshToken(token string) (AuthPayload, error)
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

func (n *noopAuth) GenerateAccessToken(payload AuthPayload) (string, error)  { return "", nil }
func (n *noopAuth) GenerateRefreshToken(payload AuthPayload) (string, error) { return "", nil }
func (n *noopAuth) ValidateAccessToken(token string) (AuthPayload, error)    { return AuthPayload{}, nil }
func (n *noopAuth) ValidateRefreshToken(token string) (AuthPayload, error)   { return AuthPayload{}, nil }
func (n *noopAuth) RevokeRefreshToken(token string) error                    { return nil }
