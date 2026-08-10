package auth

type SignupDto struct {
	LoginID   *string
	Email     *string
	Password  string
	FirstName string
	LastName  string
}

type LoginDto struct {
	LoginID    string
	Password   string
	RememberMe bool
}

type ForgotPasswordDto struct {
	Email string
}

type VerifyResetPasswordTokenDto struct {
	Token string
}

type ResetPasswordDto struct {
	Token       string
	NewPassword string
}

type UpdatePasswordDto struct {
	Id          int
	OldPassword string
	NewPassword string
}
