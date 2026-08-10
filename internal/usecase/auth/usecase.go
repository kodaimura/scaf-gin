package auth

import (
	"gorm.io/gorm"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
	"scaf-gin/internal/module"
)

type Usecase interface {
	Signup(in SignupDto) (model.Account, error)
	Login(in LoginDto) (model.Account, string, string, error)
	Refresh(refreshToken string) (core.AuthPayload, string, error)
	Logout(refreshToken string) error
	ForgotPassword(in ForgotPasswordDto) error
	VerifyResetPasswordToken(in VerifyResetPasswordTokenDto) error
	ResetPassword(in ResetPasswordDto) error
	UpdatePassword(in UpdatePasswordDto) error
}

type usecase struct {
	db                       *gorm.DB
	accountModule            module.AccountModule
	passwordResetTokenModule module.PasswordResetTokenModule
	auth                     core.AuthI
	mailer                   core.MailerI
}

func NewUsecase(
	db *gorm.DB,
	accountModule module.AccountModule,
	passwordResetTokenModule module.PasswordResetTokenModule,
	auth core.AuthI,
	mailer core.MailerI,
) Usecase {
	return &usecase{
		db:                       db,
		accountModule:            accountModule,
		passwordResetTokenModule: passwordResetTokenModule,
		auth:                     auth,
		mailer:                   mailer,
	}
}
