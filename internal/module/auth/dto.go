package auth

type SignupDto struct {
	Name     string
	Password string
}

type LoginDto struct {
	Name     string
	Password string
}

type UpdatePasswordDto struct {
	Id          int
	OldPassword string
	NewPassword string
}
