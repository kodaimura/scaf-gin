package service

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"

	"goscaf/internal/core/jwt"
	"goscaf/internal/core/logger"
	"goscaf/internal/core/errs"
	"goscaf/internal/core/utils"
	"goscaf/internal/dto"
	"goscaf/internal/model"
	"goscaf/internal/repository"
)

type AccountService interface {
	GetOne(input dto.AccountPK) (dto.Account, error)
	Delete(input dto.AccountPK) error
	Update(input dto.UpdateAccount) error
	Login(input dto.Login) (dto.Account, error)
	Signup(input dto.Signup) (dto.AccountPK, error)
	GenerateJwtPayload(input dto.AccountPK) (jwt.Payload, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService{
		accountRepository: accountRepository,
	}
}

func (srv *accountService) GetOne(input dto.AccountPK) (dto.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: input.Id})
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.Account{}, errs.NewNotFoundError()
		}
		logger.Error(err.Error())
		return dto.Account{}, errs.NewUnexpectedError(err.Error())
	}

	var ret dto.Account
	utils.MapFields(&ret, account)
	return ret, nil
}

func (srv *accountService) Update(input dto.UpdateAccount) error {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: input.Id})
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewNotFoundError()
		}
		logger.Error(err.Error())
		return errs.NewUnexpectedError(err.Error())
	}

	if input.Name != "" {
		account.Name = input.Name
	}
	if input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error(err.Error())
			return errs.NewUnexpectedError(err.Error())
		}
		account.Password = string(hashed)
	}
	if err := srv.accountRepository.Update(&account, nil); err != nil {
		if column, ok := GetConflictColumn(err); ok {
			return errs.NewConflictError(column)
		}
		logger.Error(err.Error())
		return errs.NewUnexpectedError(err.Error())
	}
	return nil
}

func (srv *accountService) Delete(input dto.AccountPK) error {
	if err := srv.accountRepository.Delete(&model.Account{Id: input.Id}, nil); err != nil {
		logger.Error(err.Error())
		return errs.NewUnexpectedError(err.Error())
	}
	return nil
}

func (srv *accountService) Login(input dto.Login) (dto.Account, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Name: input.Name})
	if err != nil {
		if err == sql.ErrNoRows {
			return dto.Account{}, errs.NewUnauthorizedError()
		}
		logger.Error(err.Error())
		return dto.Account{}, errs.NewUnexpectedError(err.Error())
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(input.Password)); err != nil {
		return dto.Account{}, errs.NewUnauthorizedError()
	}

	var ret dto.Account
	utils.MapFields(&ret, account)
	return ret, nil
}

func (srv *accountService) Signup(input dto.Signup) (dto.AccountPK, error) {
	if _, err := srv.accountRepository.GetOne(&model.Account{Name: input.Name}); err == nil {
		return dto.AccountPK{}, errs.NewConflictError("account_name")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err.Error())
		return dto.AccountPK{}, errs.NewUnexpectedError(err.Error())
	}

	account := &model.Account{
		Name:     input.Name,
		Password: string(hashed),
	}

	id, err := srv.accountRepository.Insert(account, nil)
	if err != nil {
		if column, ok := GetConflictColumn(err); ok {
			return dto.AccountPK{}, errs.NewConflictError(column)
		}
		logger.Error(err.Error())
		return dto.AccountPK{}, errs.NewUnexpectedError(err.Error())
	}

	return dto.AccountPK{Id: id}, nil
}

func (srv *accountService) GenerateJwtPayload(input dto.AccountPK) (jwt.Payload, error) {
	account, err := srv.accountRepository.GetOne(&model.Account{Id: input.Id})
	if err != nil {
		logger.Error(err.Error())
		return jwt.Payload{}, errs.NewUnexpectedError(err.Error())
	}

	cc := jwt.CustomClaims{
		AccountId:   account.Id,
		AccountName: account.Name,
	}
	return jwt.NewPayload(input.Id, int(config.JwtExpiresSeconds) ,cc), nil
}