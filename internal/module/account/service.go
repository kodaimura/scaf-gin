package account

import (
	"gorm.io/gorm"
)

type Service interface {
	Get(in GetDto, db *gorm.DB) ([]Account, error)
	GetOne(in GetOneDto, db *gorm.DB) (Account, error)
	UpdateOne(in UpdateOneDto, db *gorm.DB) (Account, error)
	DeleteOne(in DeleteOneDto, db *gorm.DB) error
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

func (srv *service) UpdateOne(in UpdateOneDto, db *gorm.DB) (Account, error) {
	acct, err := srv.GetOne(GetOneDto{Id: in.Id}, db)
	if err != nil {
		return Account{}, err
	}
	acct.Name = in.Name
	return srv.repository.Update(&acct, db)
}

func (srv *service) DeleteOne(in DeleteOneDto, db *gorm.DB) error {
	return srv.repository.Delete(&Account{Id: in.Id}, db)
}
