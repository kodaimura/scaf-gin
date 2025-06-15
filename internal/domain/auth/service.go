package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"scaf-gin/internal/core"
	"scaf-gin/internal/domain/account"
)

type Service interface {
	Signup(in SignupDto, db *gorm.DB) (account.Account, error)
	Login(in LoginDto, db *gorm.DB) (account.Account, error)
	UpdatePassword(in UpdatePasswordDto, db *gorm.DB) error
}

type service struct {
	accountRepository account.Repository
}

func NewService(accountRepository account.Repository) Service {
	return &service{
		accountRepository: accountRepository,
	}
}

func (srv *service) Signup(in SignupDto, db *gorm.DB) (account.Account, error) {
	hashed, err := hashPassword(in.Password)
	if err != nil {
		return account.Account{}, err
	}

	return srv.accountRepository.Insert(&account.Account{
		Name:     in.Name,
		Password: string(hashed),
	}, db)
}

func (srv *service) Login(in LoginDto, db *gorm.DB) (account.Account, error) {
	acct, err := srv.accountRepository.GetOne(&account.Account{Name: in.Name}, db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return account.Account{}, core.ErrUnauthorized
		}
		return account.Account{}, err
	}

	if err = verifyPassword(acct.Password, in.Password); err != nil {
		return account.Account{}, core.ErrUnauthorized
	}
	return acct, nil
}

func (srv *service) UpdatePassword(in UpdatePasswordDto, db *gorm.DB) error {
	acct, err := srv.accountRepository.GetOne(&account.Account{Id: in.Id}, db)
	if err != nil {
		return err
	}
	if err = verifyPassword(acct.Password, in.OldPassword); err != nil {
		return core.ErrBadRequest
	}

	hashed, err := hashPassword(in.NewPassword)
	if err != nil {
		return err
	}
	acct.Password = string(hashed)
	_, err = srv.accountRepository.Update(&acct, db)
	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func verifyPassword(hashedPassword, plainPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return err
	}
	return nil
}
