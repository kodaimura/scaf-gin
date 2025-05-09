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
	return srv.accountRepository.GetOne(&model.Account{Id: in.Id})
}

func (srv *accountService) CreateOne(in input.Account) (model.Account, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, err
	}

	return srv.accountRepository.Insert(&model.Account{
		Name:     in.Name,
		Password: string(hashed),
	})
}

func (srv *accountService) UpdateOne(in input.Account) (model.Account, error) {
	account, err := srv.GetOne(in)
	if err != nil {
		return model.Account{}, err
	}
	account.Name = in.Name
	return srv.accountRepository.Update(&account)
}

func (srv *accountService) DeleteOne(in input.Account) error {
	return srv.accountRepository.Delete(&model.Account{Id: in.Id})
}

func (srv *accountService) Login(in input.Login) (model.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Name: in.Name})
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return model.Account{}, core.ErrUnauthorized
		}
		return model.Account{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(in.Password)); err != nil {
		return model.Account{}, core.ErrUnauthorized
	}
	return account, nil
}

func (srv *accountService) UpdatePassword(in input.UpdatePassword) (model.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: in.Id})
	if err != nil {
		return model.Account{}, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, err
	}
	account.Password = string(hashed)
	return srv.accountRepository.Update(&account)
}
