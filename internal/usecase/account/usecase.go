package account

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"scaf-gin/internal/core"
	accountModule "scaf-gin/internal/module/account"
)

type Usecase interface {
	Get(in GetDto) ([]accountModule.Account, error)
	GetOne(in GetOneDto) (accountModule.Account, error)
	CreateOne(in CreateOneDto) (accountModule.Account, error)
	UpdateOne(in UpdateOneDto) (accountModule.Account, error)
	DisableOne(in DisableOneDto) (accountModule.Account, error)
	EnableOne(in EnableOneDto) (accountModule.Account, error)
	DeleteOne(in DeleteOneDto) error
}

type usecase struct {
	db             *gorm.DB
	accountService accountModule.Service
}

func NewUsecase(
	db *gorm.DB,
	accountService accountModule.Service,
) Usecase {
	return &usecase{
		db:             db,
		accountService: accountService,
	}
}

func (uc *usecase) Get(in GetDto) ([]accountModule.Account, error) {
	return uc.accountService.Get(accountModule.Account{}, uc.db)
}

func (uc *usecase) GetOne(in GetOneDto) (accountModule.Account, error) {
	acct, err := uc.accountService.GetOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
	if errors.Is(err, core.ErrNotFound) {
		return accountModule.Account{}, core.ErrAccountNotFound
	}
	return acct, err
}

func (uc *usecase) CreateOne(in CreateOneDto) (accountModule.Account, error) {
	loginID, err := core.ResolveLoginID(in.LoginID, in.Email)
	if err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueLoginID(loginID, 0); err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueEmail(in.Email, 0); err != nil {
		return accountModule.Account{}, err
	}

	hashed, err := hashPassword(in.Password)
	if err != nil {
		return accountModule.Account{}, err
	}

	acct, err := uc.accountService.CreateOne(accountModule.Account{
		LoginID:      loginID,
		Email:        in.Email,
		PasswordHash: hashed,
		TokenVersion: 1,
		FirstName:    in.FirstName,
		LastName:     in.LastName,
	}, uc.db)
	return acct, err
}

func (uc *usecase) UpdateOne(in UpdateOneDto) (accountModule.Account, error) {
	acct, err := uc.accountService.GetOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return accountModule.Account{}, core.ErrAccountNotFound
		}
		return accountModule.Account{}, err
	}

	loginID, err := core.ResolveLoginID(in.LoginID, in.Email)
	if err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueLoginID(loginID, acct.Id); err != nil {
		return accountModule.Account{}, err
	}
	if err := uc.ensureUniqueEmail(in.Email, acct.Id); err != nil {
		return accountModule.Account{}, err
	}

	acct.LoginID = loginID
	acct.Email = in.Email
	acct.FirstName = in.FirstName
	acct.LastName = in.LastName
	if in.Password != nil {
		hashed, err := hashPassword(*in.Password)
		if err != nil {
			return accountModule.Account{}, err
		}
		acct.PasswordHash = hashed
		acct.TokenVersion++
	}

	return uc.accountService.UpdateOne(accountModule.Account{
		Id:           acct.Id,
		LoginID:      acct.LoginID,
		Email:        acct.Email,
		PasswordHash: acct.PasswordHash,
		TokenVersion: acct.TokenVersion,
		FirstName:    acct.FirstName,
		LastName:     acct.LastName,
	}, uc.db)
}

func (uc *usecase) DisableOne(in DisableOneDto) (accountModule.Account, error) {
	acct, err := uc.accountService.DisableOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
	if errors.Is(err, core.ErrNotFound) {
		return accountModule.Account{}, core.ErrAccountNotFound
	}
	return acct, err
}

func (uc *usecase) EnableOne(in EnableOneDto) (accountModule.Account, error) {
	acct, err := uc.accountService.EnableOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
	if errors.Is(err, core.ErrNotFound) {
		return accountModule.Account{}, core.ErrAccountNotFound
	}
	return acct, err
}

func (uc *usecase) DeleteOne(in DeleteOneDto) error {
	return uc.accountService.DeleteOne(accountModule.Account{
		Id: in.Id,
	}, uc.db)
}

func (uc *usecase) ensureUniqueLoginID(loginID string, exceptID int) error {
	existing, err := uc.accountService.GetByLoginID(loginID, uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.Id != exceptID {
		return core.ErrLoginIDAlreadyExists
	}
	return nil
}

func (uc *usecase) ensureUniqueEmail(email *string, exceptID int) error {
	if email == nil || *email == "" {
		return nil
	}
	existing, err := uc.accountService.GetByEmail(*email, uc.db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return nil
		}
		return err
	}
	if existing.Id != exceptID {
		return core.ErrEmailAlreadyExists
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
