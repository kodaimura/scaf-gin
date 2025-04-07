package core

type AuthI interface {
    GenerateCredential(payload AuthPayload) (string, error)
    ValidateCredential(credential string) (AuthPayload, error)
    RevokeCredential(credential string) error
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

func (n *noopAuth) GenerateCredential(payload AuthPayload) (string, error) {
	return "", nil
}
func (n *noopAuth) ValidateCredential(credential string) (AuthPayload, error) {
	return AuthPayload{}, nil
}
func (n *noopAuth) RevokeCredential(credential string) error {
	return nil
}