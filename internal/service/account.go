package service

import (
	"errors"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"

	"goscaf/pkg/errs"
	"goscaf/internal/model"
	"goscaf/internal/repository"
	"goscaf/internal/dto/input"
)

type AccountService interface {
	GetOne(in input.Account) (model.Account, error)
	DeleteOne(in input.Account) error
	UpdateOne(in input.Account) (model.Account, error)
	Login(in input.Login) (model.Account, error)
	Signup(in input.Signup) (model.Account, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		accountRepository: accountRepository,
	}
}

func (srv *accountService) GetOne(in input.Account) (model.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{AccountId: in.AccountId})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Account{}, errs.NewNotFoundError()
		}
		return model.Account{}, errs.NewUnexpectedError(err.Error())
	}
	return account, nil
}

func (srv *accountService) UpdateOne(in input.Account) (model.Account, error) {
	account, err := srv.GetOne(in)
	if err != nil {
		return account, err
	}

	if in.AccountName != "" {
		account.AccountName = in.AccountName
	}
	if in.AccountPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(in.AccountPassword), bcrypt.DefaultCost)
		if err != nil {
			return model.Account{}, errs.NewUnexpectedError(err.Error())
		}
		account.AccountPassword = string(hashed)
	}
	account, err = srv.accountRepository.Update(&account)
	if  err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.Account{}, errs.NewConflictError()
		}
		return model.Account{}, errs.NewUnexpectedError(err.Error())
	}
	return account, nil
}

func (srv *accountService) DeleteOne(in input.Account) error {
	if err := srv.accountRepository.Delete(&model.Account{AccountId: in.AccountId}); err != nil {
		return errs.NewUnexpectedError(err.Error())
	}
	return nil
}

func (srv *accountService) Login(in input.Login) (model.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{AccountName: in.AccountName})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Account{}, errs.NewUnauthorizedError()
		}
		return model.Account{}, errs.NewUnexpectedError(err.Error())
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.AccountPassword), []byte(in.AccountPassword)); err != nil {
		return model.Account{}, errs.NewUnauthorizedError()
	}
	return account, nil
}

func (srv *accountService) Signup(in input.Signup) (model.Account, error) {
	if _, err := srv.accountRepository.GetOne(&model.Account{AccountName: in.AccountName}); err == nil {
		return model.Account{}, errs.NewConflictError()
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(in.AccountPassword), bcrypt.DefaultCost)
	if err != nil {
		return model.Account{}, errs.NewUnexpectedError(err.Error())
	}

	account := model.Account{
		AccountName:     in.AccountName,
		AccountPassword: string(hashed),
	}

	account, err = srv.accountRepository.Insert(&account)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return model.Account{}, errs.NewConflictError()
		}
		return model.Account{}, errs.NewUnexpectedError(err.Error())
	}

	return account, nil
}