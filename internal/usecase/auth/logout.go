package auth

func (uc *usecase) Logout(refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return uc.auth.RevokeRefreshToken(refreshToken)
}
