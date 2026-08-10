package module

import (
	"time"

	"gorm.io/gorm"

	"scaf-gin/internal/helper"
	"scaf-gin/internal/model"
)

type PasswordResetTokenModule interface {
	Create(entity model.PasswordResetToken) (model.PasswordResetToken, error)
	GetByHash(tokenHash string) (model.PasswordResetToken, error)
	FindLatestByAccountID(accountID int) (model.PasswordResetToken, error)
	InvalidateActiveTokens(accountID int, now time.Time) error
	Update(entity model.PasswordResetToken) (model.PasswordResetToken, error)
	WithTx(tx *gorm.DB) PasswordResetTokenModule
}

type passwordResetTokenModule struct {
	db *gorm.DB
}

func NewPasswordResetTokenModule(db *gorm.DB) PasswordResetTokenModule {
	return &passwordResetTokenModule{db: db}
}

func (m *passwordResetTokenModule) WithTx(tx *gorm.DB) PasswordResetTokenModule {
	return &passwordResetTokenModule{db: tx}
}

func (m *passwordResetTokenModule) Create(entity model.PasswordResetToken) (model.PasswordResetToken, error) {
	err := m.db.Create(&entity).Error
	return entity, helper.HandleGormError(err)
}

func (m *passwordResetTokenModule) GetByHash(tokenHash string) (model.PasswordResetToken, error) {
	var token model.PasswordResetToken
	err := m.db.Where("token_hash = ?", tokenHash).First(&token).Error
	return token, helper.HandleGormError(err)
}

func (m *passwordResetTokenModule) FindLatestByAccountID(accountID int) (model.PasswordResetToken, error) {
	var token model.PasswordResetToken
	err := m.db.Where("account_id = ?", accountID).Order("created_at DESC").First(&token).Error
	return token, helper.HandleGormError(err)
}

func (m *passwordResetTokenModule) InvalidateActiveTokens(accountID int, now time.Time) error {
	err := m.db.Model(&model.PasswordResetToken{}).
		Where("account_id = ?", accountID).
		Where("used_at IS NULL").
		Where("expires_at > ?", now).
		Update("used_at", now).
		Error
	return helper.HandleGormError(err)
}

func (m *passwordResetTokenModule) Update(entity model.PasswordResetToken) (model.PasswordResetToken, error) {
	err := m.db.Save(&entity).Error
	return entity, helper.HandleGormError(err)
}
