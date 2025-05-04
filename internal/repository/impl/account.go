package impl

import (
	"gorm.io/gorm"

	"scaf-gin/internal/model"
)

type gormAccountRepository struct {
	db *gorm.DB
}

func NewGormAccountRepository(db *gorm.DB) *gormAccountRepository {
	return &gormAccountRepository{db: db}
}

func (rep *gormAccountRepository) Get(a *model.Account) ([]model.Account, error) {
	var accounts []model.Account
	err := rep.db.Find(&accounts, a).Error
	return accounts, handleGormError(err)
}

func (rep *gormAccountRepository) GetOne(a *model.Account) (model.Account, error) {
	var account model.Account
	err := rep.db.First(&account, a).Error
	return account, handleGormError(err)
}

func (rep *gormAccountRepository) GetAll(a *model.Account) ([]model.Account, error) {
	var accounts []model.Account
	err := rep.db.Unscoped().Find(&accounts, a).Error
	return accounts, handleGormError(err)
}

func (rep *gormAccountRepository) Insert(a *model.Account) (model.Account, error) {
	err := rep.db.Create(a).Error
	return *a, handleGormError(err)
}

func (rep *gormAccountRepository) Update(a *model.Account) (model.Account, error) {
	err := rep.db.Save(a).Error
	return *a, handleGormError(err)
}

func (rep *gormAccountRepository) Delete(a *model.Account) error {
	err := rep.db.Delete(a).Error
	return handleGormError(err)
}

func (rep *gormAccountRepository) InsertTx(a *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Create(a).Error
	return *a, handleGormError(err)
}

func (rep *gormAccountRepository) UpdateTx(a *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Save(a).Error
	return *a, handleGormError(err)
}

func (rep *gormAccountRepository) DeleteTx(a *model.Account, tx *gorm.DB) error {
	err := tx.Delete(a).Error
	return handleGormError(err)
}
