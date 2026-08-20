package auth

import (
	"gorm.io/gorm"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
	"scaf-gin/internal/module"
)

type Usecase interface {
	AuthorizeAccessToken(payload core.AuthPayload) (int64, error)
	Signup(in SignupDto) (model.Account, error)
	Login(in LoginDto) (model.Account, string, string, error)
	Refresh(refreshToken string) (core.AuthPayload, string, error)
	ForgotPassword(in ForgotPasswordDto) error
	VerifyResetPasswordToken(in VerifyResetPasswordTokenDto) error
	ResetPassword(in ResetPasswordDto) error
	UpdatePassword(in UpdatePasswordDto) error
}

type usecase struct {
	db                       *gorm.DB
	accountModule            module.AccountModule
	passwordResetTokenModule module.PasswordResetTokenModule
	auth                     core.Auth
	mailer                   core.Mailer
}

func NewUsecase(
	db *gorm.DB,
	accountModule module.AccountModule,
	passwordResetTokenModule module.PasswordResetTokenModule,
	auth core.Auth,
	mailer core.Mailer,
) Usecase {
	return &usecase{
		db:                       db,
		accountModule:            accountModule,
		passwordResetTokenModule: passwordResetTokenModule,
		auth:                     auth,
		mailer:                   mailer,
	}
}
