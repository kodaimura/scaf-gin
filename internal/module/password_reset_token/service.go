package password_reset_token

import (
	"time"

	"gorm.io/gorm"
)

type Service interface {
	GetByHash(tokenHash string, db *gorm.DB) (PasswordResetToken, error)
	FindLatestByAccountId(accountId int, db *gorm.DB) (PasswordResetToken, error)
	CreateOne(in PasswordResetToken, db *gorm.DB) (PasswordResetToken, error)
	UpdateOne(in PasswordResetToken, db *gorm.DB) (PasswordResetToken, error)
	InvalidateActiveTokens(accountId int, now time.Time, db *gorm.DB) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (srv *service) GetByHash(tokenHash string, db *gorm.DB) (PasswordResetToken, error) {
	return srv.repository.GetByHash(tokenHash, db)
}

func (srv *service) FindLatestByAccountId(accountId int, db *gorm.DB) (PasswordResetToken, error) {
	return srv.repository.FindLatestByAccountId(accountId, db)
}

func (srv *service) CreateOne(in PasswordResetToken, db *gorm.DB) (PasswordResetToken, error) {
	return srv.repository.Insert(&in, db)
}

func (srv *service) UpdateOne(in PasswordResetToken, db *gorm.DB) (PasswordResetToken, error) {
	return srv.repository.Update(&in, db)
}

func (srv *service) InvalidateActiveTokens(accountId int, now time.Time, db *gorm.DB) error {
	return srv.repository.InvalidateActiveTokens(accountId, now, db)
}
