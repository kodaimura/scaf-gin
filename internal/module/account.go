package module

import (
	"time"

	"gorm.io/gorm"

	"scaf-gin/internal/core"
	"scaf-gin/internal/model"
)

type AccountModule interface {
	Create(entity model.Account) (model.Account, error)
	GetAll() ([]model.Account, error)
	GetByID(accountID int) (model.Account, error)
	GetByEmail(email string) (model.Account, error)
	GetByLoginID(loginID string) (model.Account, error)
	Update(entity model.Account) (model.Account, error)
	Disable(entity model.Account) (model.Account, error)
	Enable(entity model.Account) (model.Account, error)
	Delete(entity model.Account) error
	WithTx(tx *gorm.DB) AccountModule
}

type accountModule struct {
	db *gorm.DB
}

func NewAccountModule(db *gorm.DB) AccountModule {
	return &accountModule{db: db}
}

func (m *accountModule) WithTx(tx *gorm.DB) AccountModule {
	return &accountModule{db: tx}
}

func (m *accountModule) Create(entity model.Account) (model.Account, error) {
	if entity.TokenVersion == 0 {
		entity.TokenVersion = 1
	}
	err := m.db.Create(&entity).Error
	return entity, core.HandleGormError(err)
}

func (m *accountModule) GetAll() ([]model.Account, error) {
	var accounts []model.Account
	err := m.db.Order("id").Find(&accounts).Error
	return accounts, core.HandleGormError(err)
}

func (m *accountModule) GetByID(accountID int) (model.Account, error) {
	var account model.Account
	err := m.db.First(&account, model.Account{ID: accountID}).Error
	return account, core.HandleGormError(err)
}

func (m *accountModule) GetByEmail(email string) (model.Account, error) {
	var account model.Account
	err := m.db.Where("email = ?", email).First(&account).Error
	return account, core.HandleGormError(err)
}

func (m *accountModule) GetByLoginID(loginID string) (model.Account, error) {
	var account model.Account
	err := m.db.Where("login_id = ?", loginID).First(&account).Error
	return account, core.HandleGormError(err)
}

func (m *accountModule) Update(entity model.Account) (model.Account, error) {
	err := m.db.Save(&entity).Error
	return entity, core.HandleGormError(err)
}

func (m *accountModule) Disable(entity model.Account) (model.Account, error) {
	now := time.Now()
	entity.DisabledAt = &now
	entity.TokenVersion++
	return m.Update(entity)
}

func (m *accountModule) Enable(entity model.Account) (model.Account, error) {
	entity.DisabledAt = nil
	return m.Update(entity)
}

func (m *accountModule) Delete(entity model.Account) error {
	err := m.db.Delete(&entity).Error
	return core.HandleGormError(err)
}
