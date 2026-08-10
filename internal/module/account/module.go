package account

import (
	"time"

	"gorm.io/gorm"

	"scaf-gin/internal/helper"
)

type Module interface {
	Create(entity Account) (Account, error)
	GetAll() ([]Account, error)
	GetByID(accountID int) (Account, error)
	GetByEmail(email string) (Account, error)
	GetByLoginID(loginID string) (Account, error)
	Update(entity Account) (Account, error)
	Disable(entity Account) (Account, error)
	Enable(entity Account) (Account, error)
	Delete(entity Account) error
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

func (m *module) Create(entity Account) (Account, error) {
	if entity.TokenVersion == 0 {
		entity.TokenVersion = 1
	}
	err := m.db.Create(&entity).Error
	return entity, helper.HandleGormError(err)
}

func (m *module) GetAll() ([]Account, error) {
	var accounts []Account
	err := m.db.Order("id").Find(&accounts).Error
	return accounts, helper.HandleGormError(err)
}

func (m *module) GetByID(accountID int) (Account, error) {
	var account Account
	err := m.db.First(&account, Account{Id: accountID}).Error
	return account, helper.HandleGormError(err)
}

func (m *module) GetByEmail(email string) (Account, error) {
	var account Account
	err := m.db.Where("email = ?", email).First(&account).Error
	return account, helper.HandleGormError(err)
}

func (m *module) GetByLoginID(loginID string) (Account, error) {
	var account Account
	err := m.db.Where("login_id = ?", loginID).First(&account).Error
	return account, helper.HandleGormError(err)
}

func (m *module) Update(entity Account) (Account, error) {
	err := m.db.Save(&entity).Error
	return entity, helper.HandleGormError(err)
}

func (m *module) Disable(entity Account) (Account, error) {
	now := time.Now()
	entity.DisabledAt = &now
	entity.TokenVersion++
	return m.Update(entity)
}

func (m *module) Enable(entity Account) (Account, error) {
	entity.DisabledAt = nil
	return m.Update(entity)
}

func (m *module) Delete(entity Account) error {
	err := m.db.Delete(&entity).Error
	return helper.HandleGormError(err)
}
