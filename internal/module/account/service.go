package account

import (
	"time"

	"gorm.io/gorm"
)

type Service interface {
	Get(in Account, db *gorm.DB) ([]Account, error)
	GetOne(in Account, db *gorm.DB) (Account, error)
	GetByEmail(email string, db *gorm.DB) (Account, error)
	GetByLoginID(loginID string, db *gorm.DB) (Account, error)
	CreateOne(in Account, db *gorm.DB) (Account, error)
	UpdateOne(in Account, db *gorm.DB) (Account, error)
	DisableOne(in Account, db *gorm.DB) (Account, error)
	EnableOne(in Account, db *gorm.DB) (Account, error)
	DeleteOne(in Account, db *gorm.DB) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (srv *service) Get(in Account, db *gorm.DB) ([]Account, error) {
	return srv.repository.GetAll(&in, db)
}

func (srv *service) GetOne(in Account, db *gorm.DB) (Account, error) {
	return srv.repository.GetOne(&in, db)
}

func (srv *service) GetByEmail(email string, db *gorm.DB) (Account, error) {
	return srv.repository.GetByEmail(email, db)
}

func (srv *service) GetByLoginID(loginID string, db *gorm.DB) (Account, error) {
	return srv.repository.GetByLoginID(loginID, db)
}

func (srv *service) CreateOne(in Account, db *gorm.DB) (Account, error) {
	return srv.repository.Insert(&in, db)
}

func (srv *service) UpdateOne(in Account, db *gorm.DB) (Account, error) {
	acct, err := srv.repository.GetOne(&Account{Id: in.Id}, db)
	if err != nil {
		return Account{}, err
	}
	acct.Email = in.Email
	acct.LoginID = in.LoginID
	acct.FirstName = in.FirstName
	acct.LastName = in.LastName
	if in.PasswordHash != "" {
		acct.PasswordHash = in.PasswordHash
	}
	if in.TokenVersion > 0 {
		acct.TokenVersion = in.TokenVersion
	}
	return srv.repository.Update(&acct, db)
}

func (srv *service) DisableOne(in Account, db *gorm.DB) (Account, error) {
	acct, err := srv.repository.GetOne(&Account{Id: in.Id}, db)
	if err != nil {
		return Account{}, err
	}
	now := time.Now()
	acct.DisabledAt = &now
	acct.TokenVersion++
	return srv.repository.Update(&acct, db)
}

func (srv *service) EnableOne(in Account, db *gorm.DB) (Account, error) {
	acct, err := srv.repository.GetOne(&Account{Id: in.Id}, db)
	if err != nil {
		return Account{}, err
	}
	acct.DisabledAt = nil
	return srv.repository.Update(&acct, db)
}

func (srv *service) DeleteOne(in Account, db *gorm.DB) error {
	return srv.repository.Delete(&Account{Id: in.Id}, db)
}
