package account

import (
	"gorm.io/gorm"

	"scaf-gin/internal/helper"
)

type Repository interface {
	Get(m *Account, db *gorm.DB) ([]Account, error)
	GetOne(m *Account, db *gorm.DB) (Account, error)
	GetAll(m *Account, db *gorm.DB) ([]Account, error)
	Insert(m *Account, db *gorm.DB) (Account, error)
	Update(m *Account, db *gorm.DB) (Account, error)
	Delete(m *Account, db *gorm.DB) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (rep *repository) Get(m *Account, db *gorm.DB) ([]Account, error) {
	var accounts []Account
	err := db.Find(&accounts, m).Error
	return accounts, helper.HandleGormError(err)
}

func (rep *repository) GetOne(m *Account, db *gorm.DB) (Account, error) {
	var account Account
	err := db.First(&account, m).Error
	return account, helper.HandleGormError(err)
}

func (rep *repository) GetAll(m *Account, db *gorm.DB) ([]Account, error) {
	var accounts []Account
	err := db.Unscoped().Find(&accounts, m).Error
	return accounts, helper.HandleGormError(err)
}

func (rep *repository) Insert(m *Account, db *gorm.DB) (Account, error) {
	err := db.Create(m).Error
	return *m, helper.HandleGormError(err)
}

func (rep *repository) Update(m *Account, db *gorm.DB) (Account, error) {
	err := db.Save(m).Error
	return *m, helper.HandleGormError(err)
}

func (rep *repository) Delete(m *Account, db *gorm.DB) error {
	err := db.Delete(m).Error
	return helper.HandleGormError(err)
}
