package password_reset_token

import (
	"time"

	"gorm.io/gorm"

	"scaf-gin/internal/helper"
)

type Module interface {
	Create(entity PasswordResetToken) (PasswordResetToken, error)
	GetByHash(tokenHash string) (PasswordResetToken, error)
	FindLatestByAccountID(accountID int) (PasswordResetToken, error)
	InvalidateActiveTokens(accountID int, now time.Time) error
	Update(entity PasswordResetToken) (PasswordResetToken, error)
	WithTx(tx *gorm.DB) Module
}

type module struct {
	db *gorm.DB
}

func NewModule(db *gorm.DB) Module {
	return &module{db: db}
}

func (m *module) WithTx(tx *gorm.DB) Module {
	return &module{db: tx}
}

func (m *module) Create(entity PasswordResetToken) (PasswordResetToken, error) {
	err := m.db.Create(&entity).Error
	return entity, helper.HandleGormError(err)
}

func (m *module) GetByHash(tokenHash string) (PasswordResetToken, error) {
	var token PasswordResetToken
	err := m.db.Where("token_hash = ?", tokenHash).First(&token).Error
	return token, helper.HandleGormError(err)
}

func (m *module) FindLatestByAccountID(accountID int) (PasswordResetToken, error) {
	var token PasswordResetToken
	err := m.db.Where("account_id = ?", accountID).Order("created_at DESC").First(&token).Error
	return token, helper.HandleGormError(err)
}

func (m *module) InvalidateActiveTokens(accountID int, now time.Time) error {
	err := m.db.Model(&PasswordResetToken{}).
		Where("account_id = ?", accountID).
		Where("used_at IS NULL").
		Where("expires_at > ?", now).
		Update("used_at", now).
		Error
	return helper.HandleGormError(err)
}

func (m *module) Update(entity PasswordResetToken) (PasswordResetToken, error) {
	err := m.db.Save(&entity).Error
	return entity, helper.HandleGormError(err)
}
