package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"scaf-gin/internal/core"
	"scaf-gin/internal/dto/input"
	"scaf-gin/internal/model"
	"scaf-gin/internal/repository"
)

type AccountService interface {
	Get(in input.Account) ([]model.Account, error)
	GetOne(in input.Account) (model.Account, error)
	CreateOne(in input.Account) (model.Account, error)
	UpdateOne(in input.Account) (model.Account, error)
	DeleteOne(in input.Account) error

	Login(in input.Login) (model.Account, error)
	UpdatePassword(in input.UpdatePassword) (model.Account, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		accountRepository: accountRepository,
	}
}

func (srv *accountService) Get(in input.Account) ([]model.Account, error) {
	return srv.accountRepository.Get(&model.Account{})
}

func (srv *accountService) GetOne(in input.Account) (model.Account, error) {
	return srv.accountRepository.GetOne(&model.Account{AccountId: in.AccountId})
}

func (srv *accountService) CreateOne(in input.Account) (model.Account, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.AccountPassword), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, err
	}

	return srv.accountRepository.Insert(&model.Account{
		AccountName:     in.AccountName,
		AccountPassword: string(hashed),
	})
}

func (srv *accountService) UpdateOne(in input.Account) (model.Account, error) {
	account, err := srv.GetOne(in)
	if err != nil {
		return model.Account{}, err
	}
	account.AccountName = in.AccountName
	return srv.accountRepository.Update(&account)
}

func (srv *accountService) DeleteOne(in input.Account) error {
	return srv.accountRepository.Delete(&model.Account{AccountId: in.AccountId})
}

func (srv *accountService) Login(in input.Login) (model.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{AccountName: in.AccountName})
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return model.Account{}, core.ErrUnauthorized
		}
		return model.Account{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.AccountPassword), []byte(in.AccountPassword)); err != nil {
		return model.Account{}, core.ErrUnauthorized
	}
	return account, nil
}

func (srv *accountService) UpdatePassword(in input.UpdatePassword) (model.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{AccountId: in.AccountId})
	if err != nil {
		return model.Account{}, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.AccountPassword), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, err
	}
	account.AccountPassword = string(hashed)
	return srv.accountRepository.Update(&account)
}
