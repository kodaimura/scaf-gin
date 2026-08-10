package password_reset_token

import (
	"time"

	"gorm.io/gorm"

	"scaf-gin/internal/helper"
)

type Repository interface {
	GetByHash(tokenHash string, db *gorm.DB) (PasswordResetToken, error)
	FindLatestByAccountId(accountId int, db *gorm.DB) (PasswordResetToken, error)
	Insert(m *PasswordResetToken, db *gorm.DB) (PasswordResetToken, error)
	Update(m *PasswordResetToken, db *gorm.DB) (PasswordResetToken, error)
	InvalidateActiveTokens(accountId int, now time.Time, db *gorm.DB) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (rep *repository) GetByHash(tokenHash string, db *gorm.DB) (PasswordResetToken, error) {
	var token PasswordResetToken
	err := db.Where("token_hash = ?", tokenHash).First(&token).Error
	return token, helper.HandleGormError(err)
}

func (rep *repository) FindLatestByAccountId(accountId int, db *gorm.DB) (PasswordResetToken, error) {
	var token PasswordResetToken
	err := db.Where("account_id = ?", accountId).Order("created_at DESC").First(&token).Error
	return token, helper.HandleGormError(err)
}

func (rep *repository) Insert(m *PasswordResetToken, db *gorm.DB) (PasswordResetToken, error) {
	err := db.Create(m).Error
	return *m, helper.HandleGormError(err)
}

func (rep *repository) Update(m *PasswordResetToken, db *gorm.DB) (PasswordResetToken, error) {
	err := db.Save(m).Error
	return *m, helper.HandleGormError(err)
}

func (rep *repository) InvalidateActiveTokens(accountId int, now time.Time, db *gorm.DB) error {
	err := db.Model(&PasswordResetToken{}).
		Where("account_id = ?", accountId).
		Where("used_at IS NULL").
		Where("expires_at > ?", now).
		Update("used_at", now).
		Error
	return helper.HandleGormError(err)
}
