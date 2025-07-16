package account_profile

import (
	"gorm.io/gorm"

	"scaf-gin/internal/helper"
)

type Repository interface {
	Get(m *AccountProfile, db *gorm.DB) ([]AccountProfile, error)
	GetOne(m *AccountProfile, db *gorm.DB) (AccountProfile, error)
	GetAll(m *AccountProfile, db *gorm.DB) ([]AccountProfile, error)
	Insert(m *AccountProfile, db *gorm.DB) (AccountProfile, error)
	Update(m *AccountProfile, db *gorm.DB) (AccountProfile, error)
	Delete(m *AccountProfile, db *gorm.DB) error
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}

func (rep *repository) Get(m *AccountProfile, db *gorm.DB) ([]AccountProfile, error) {
	var accounts []AccountProfile
	err := db.Find(&accounts, m).Error
	return accounts, helper.HandleGormError(err)
}

func (rep *repository) GetOne(m *AccountProfile, db *gorm.DB) (AccountProfile, error) {
	var account AccountProfile
	err := db.First(&account, m).Error
	return account, helper.HandleGormError(err)
}

func (rep *repository) GetAll(m *AccountProfile, db *gorm.DB) ([]AccountProfile, error) {
	var accounts []AccountProfile
	err := db.Unscoped().Find(&accounts, m).Error
	return accounts, helper.HandleGormError(err)
}

func (rep *repository) Insert(m *AccountProfile, db *gorm.DB) (AccountProfile, error) {
	err := db.Create(m).Error
	return *m, helper.HandleGormError(err)
}

func (rep *repository) Update(m *AccountProfile, db *gorm.DB) (AccountProfile, error) {
	err := db.Save(m).Error
	return *m, helper.HandleGormError(err)
}

func (rep *repository) Delete(m *AccountProfile, db *gorm.DB) error {
	err := db.Delete(m).Error
	return helper.HandleGormError(err)
}
