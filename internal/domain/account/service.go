package account

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"scaf-gin/internal/core"
)

type Service interface {
	Get(in GetDto, db *gorm.DB) ([]Account, error)
	GetOne(in GetOneDto, db *gorm.DB) (Account, error)
	CreateOne(in CreateOneDto, db *gorm.DB) (Account, error)
	UpdateOne(in UpdateOneDto, db *gorm.DB) (Account, error)
	DeleteOne(in DeleteOneDto, db *gorm.DB) error

	Login(in LoginDto, db *gorm.DB) (Account, error)
	UpdatePassword(in UpdatePasswordDto, db *gorm.DB) (Account, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (srv *service) Get(in GetDto, db *gorm.DB) ([]Account, error) {
	return srv.repository.Get(&Account{}, db)
}

func (srv *service) GetOne(in GetOneDto, db *gorm.DB) (Account, error) {
	return srv.repository.GetOne(&Account{Id: in.Id}, db)
}

func (srv *service) CreateOne(in CreateOneDto, db *gorm.DB) (Account, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}

	return srv.repository.Insert(&Account{
		Name:     in.Name,
		Password: string(hashed),
	}, db)
}

func (srv *service) UpdateOne(in UpdateOneDto, db *gorm.DB) (Account, error) {
	account, err := srv.GetOne(GetOneDto{Id: in.Id}, db)
	if err != nil {
		return Account{}, err
	}
	account.Name = in.Name
	return srv.repository.Update(&account, db)
}

func (srv *service) DeleteOne(in DeleteOneDto, db *gorm.DB) error {
	return srv.repository.Delete(&Account{Id: in.Id}, db)
}

func (srv *service) Login(in LoginDto, db *gorm.DB) (Account, error) {
	account, err := srv.repository.GetOne(&Account{Name: in.Name}, db)
	if err != nil {
		if errors.Is(err, core.ErrNotFound) {
			return Account{}, core.ErrUnauthorized
		}
		return Account{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(in.Password)); err != nil {
		return Account{}, core.ErrUnauthorized
	}
	return account, nil
}

func (srv *service) UpdatePassword(in UpdatePasswordDto, db *gorm.DB) (Account, error) {
	account, err := srv.repository.GetOne(&Account{Id: in.Id}, db)
	if err != nil {
		return Account{}, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return Account{}, err
	}
	account.Password = string(hashed)
	return srv.repository.Update(&account, db)
}
