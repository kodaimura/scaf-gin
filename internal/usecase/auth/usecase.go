package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
	profileModule "scaf-gin/internal/module/account_profile"
)

type Usecase interface {
	Signup(in SignupDto) (accountModule.Account, error)
	Login(in LoginDto) (accountModule.Account, string, string, error)
	Refresh(refreshToken string) (core.AuthPayload, string, error)
	UpdatePassword(in UpdatePasswordDto) error
}

type usecase struct {
	db                    *gorm.DB
	accountService        accountModule.Service
	accountProfileService profileModule.Service
}

func NewUsecase(
	db *gorm.DB,
	accountService accountModule.Service,
	accountProfileService profileModule.Service,
) Usecase {
	return &usecase{
		db:                    db,
		accountService:        accountService,
		accountProfileService: accountProfileService,
	}
}

func (uc *usecase) Signup(in SignupDto) (accountModule.Account, error) {
	var account accountModule.Account
	err := uc.db.Transaction(func(tx *gorm.DB) error {
		hashed, err := hashPassword(in.Password)
		if err != nil {
			return err
		}

		account, err = uc.accountService.CreateOne(accountModule.Account{
			Name:     in.Name,
			Password: string(hashed),
		}, tx)

		if err != nil {
			return err
		}

		_, err = uc.accountProfileService.CreateOne(profileModule.AccountProfile{
			AccountId:   account.Id,
			DisplayName: account.Name,
			Bio:         "",
			AvatarURL:   "",
		}, tx)
		return err
	})

	return account, err
}

func (uc *usecase) Login(in LoginDto) (accountModule.Account, string, string, error) {
	acct, err := uc.accountService.GetOne(accountModule.Account{
		Name: in.Name,
	}, uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return accountModule.Account{}, "", "", core.ErrUnauthorized
		}
		return accountModule.Account{}, "", "", err
	}

	if !verifyPassword(acct.Password, in.Password) {
		return accountModule.Account{}, "", "", core.ErrUnauthorized
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:   acct.Id,
		AccountName: acct.Name,
	})
	if err != nil {
		return accountModule.Account{}, "", "", err
	}

	refreshToken, err := core.Auth.CreateRefreshToken(core.AuthPayload{
		AccountId:   acct.Id,
		AccountName: acct.Name,
	})
	if err != nil {
		return accountModule.Account{}, "", "", err
	}
	return acct, accessToken, refreshToken, nil
}

func (uc *usecase) Refresh(refreshToken string) (core.AuthPayload, string, error) {
	payload, err := core.Auth.VerifyRefreshToken(refreshToken)
	if err != nil {
		return core.AuthPayload{}, "", core.ErrUnauthorized
	}

	accessToken, err := core.Auth.CreateAccessToken(core.AuthPayload{
		AccountId:   payload.AccountId,
		AccountName: payload.AccountName,
	})

	return payload, accessToken, err
}

func (uc *usecase) UpdatePassword(in UpdatePasswordDto) error {
	acct, err := uc.accountService.GetOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
	if err != nil {
		return err
	}
	if !verifyPassword(acct.Password, in.OldPassword) {
		return core.ErrBadRequest
	}

	hashed, err := hashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	acct.Password = string(hashed)
	_, err = uc.accountService.UpdateOne(acct, uc.db)
	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(hashedPassword, plainPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return false
	}
	return true
}
