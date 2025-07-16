package account

import (
	"gorm.io/gorm"
)

type Service interface {
	Get(in Account, db *gorm.DB) ([]Account, error)
	GetOne(in Account, db *gorm.DB) (Account, error)
	CreateOne(in Account, db *gorm.DB) (Account, error)
	UpdateOne(in Account, db *gorm.DB) (Account, error)
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
	return srv.repository.Get(&Account{}, db)
}

func (srv *service) GetOne(in Account, db *gorm.DB) (Account, error) {
	return srv.repository.GetOne(&in, db)
}

func (srv *service) CreateOne(in Account, db *gorm.DB) (Account, error) {
	return srv.repository.Insert(&in, db)
}

func (srv *service) UpdateOne(in Account, db *gorm.DB) (Account, error) {
	acct, err := srv.repository.GetOne(&Account{Id: in.Id}, db)
	if err != nil {
		return Account{}, err
	}
	acct.Name = in.Name
	return srv.repository.Update(&acct, db)
}

func (srv *service) DeleteOne(in Account, db *gorm.DB) error {
	return srv.repository.Delete(&Account{Id: in.Id}, db)
}
