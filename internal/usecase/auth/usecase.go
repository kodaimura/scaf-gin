package auth

import (
	"gorm.io/gorm"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
	passwordResetTokenModule "scaf-gin/internal/module/password_reset_token"
)

type Usecase interface {
	Signup(in SignupDto) (accountModule.Account, error)
	Login(in LoginDto) (accountModule.Account, string, string, error)
	Refresh(refreshToken string) (core.AuthPayload, string, error)
	ForgotPassword(in ForgotPasswordDto) error
	VerifyResetPasswordToken(in VerifyResetPasswordTokenDto) error
	ResetPassword(in ResetPasswordDto) error
	UpdatePassword(in UpdatePasswordDto) error
}

type usecase struct {
	db                       *gorm.DB
	accountModule            accountModule.Module
	passwordResetTokenModule passwordResetTokenModule.Module
}

func NewUsecase(
	db *gorm.DB,
	accountModule accountModule.Module,
	passwordResetTokenModule passwordResetTokenModule.Module,
) Usecase {
	return &usecase{
		db:                       db,
		accountModule:            accountModule,
		passwordResetTokenModule: passwordResetTokenModule,
	}
}
