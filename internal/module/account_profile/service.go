package account_profile

import (
	"gorm.io/gorm"
)

type Service interface {
	Get(in AccountProfile, db *gorm.DB) ([]AccountProfile, error)
	GetOne(in AccountProfile, db *gorm.DB) (AccountProfile, error)
	CreateOne(in AccountProfile, db *gorm.DB) (AccountProfile, error)
	UpdateOne(in AccountProfile, db *gorm.DB) (AccountProfile, error)
	DeleteOne(in AccountProfile, db *gorm.DB) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (srv *service) Get(in AccountProfile, db *gorm.DB) ([]AccountProfile, error) {
	return srv.repository.Get(&AccountProfile{}, db)
}

func (srv *service) GetOne(in AccountProfile, db *gorm.DB) (AccountProfile, error) {
	return srv.repository.GetOne(&AccountProfile{
		AccountId: in.AccountId,
	}, db)
}

func (srv *service) CreateOne(in AccountProfile, db *gorm.DB) (AccountProfile, error) {
	return srv.repository.Insert(&AccountProfile{
		AccountId:   in.AccountId,
		DisplayName: in.DisplayName,
		Bio:         in.Bio,
		AvatarURL:   in.AvatarURL,
	}, db)
}

func (srv *service) UpdateOne(in AccountProfile, db *gorm.DB) (AccountProfile, error) {
	profile, err := srv.repository.GetOne(&AccountProfile{
		AccountId: in.AccountId,
	}, db)
	if err != nil {
		return AccountProfile{}, err
	}

	profile.DisplayName = in.DisplayName
	profile.Bio = in.Bio
	profile.AvatarURL = in.AvatarURL

	return srv.repository.Update(&profile, db)
}

func (srv *service) DeleteOne(in AccountProfile, db *gorm.DB) error {
	return srv.repository.Delete(&AccountProfile{
		AccountId: in.AccountId,
	}, db)
}
