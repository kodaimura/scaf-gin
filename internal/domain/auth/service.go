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
	UpdatePassword(in UpdatePasswordDto, db *gorm.DB) (account.Account, error)
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
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return account.Account{}, err
	}

	return srv.accountRepository.Insert(&account.Account{
		Name:     in.Name,
		Password: string(hashed),
	}, db)
}

func (srv *service) Login(in LoginDto, db *gorm.DB) (account.Account, error) {
	a, err := srv.accountRepository.GetOne(&account.Account{Name: in.Name}, db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return account.Account{}, core.ErrUnauthorized
		}
		return account.Account{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(in.Password)); err != nil {
		return account.Account{}, core.ErrUnauthorized
	}
	return a, nil
}

func (srv *service) UpdatePassword(in UpdatePasswordDto, db *gorm.DB) (account.Account, error) {
	a, err := srv.accountRepository.GetOne(&account.Account{Id: in.Id}, db)
	if err != nil {
		return account.Account{}, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return account.Account{}, err
	}
	a.Password = string(hashed)
	return srv.accountRepository.Update(&a, db)
}
