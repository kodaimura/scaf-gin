package core

type AuthI interface {
    GenerateToken(payload AuthPayload) (string, error)
    ValidateToken(token string) (AuthPayload, error)
    RevokeToken(token string) error
}

type AuthPayload struct {
    AccountId int
	AccountName string
}

var Auth AuthI = &noopAuth{}

func SetAuth(l AuthI) {
	Auth = l
}

type noopAuth struct {}

func (n *noopAuth) GenerateToken(payload AuthPayload) (string, error) {
	return "", nil
}
func (n *noopAuth) ValidateToken(token string) (AuthPayload, error) {
	return AuthPayload{}, nil
}
func (n *noopAuth) RevokeToken(token string) error {
	return nil
}