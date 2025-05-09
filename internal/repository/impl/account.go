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

func (rep *gormAccountRepository) Get(m *model.Account) ([]model.Account, error) {
	var accounts []model.Account
	err := rep.db.Find(&accounts, m).Error
	return accounts, handleGormError(err)
}

func (rep *gormAccountRepository) GetOne(m *model.Account) (model.Account, error) {
	var account model.Account
	err := rep.db.First(&account, m).Error
	return account, handleGormError(err)
}

func (rep *gormAccountRepository) GetAll(m *model.Account) ([]model.Account, error) {
	var accounts []model.Account
	err := rep.db.Unscoped().Find(&accounts, m).Error
	return accounts, handleGormError(err)
}

func (rep *gormAccountRepository) Insert(m *model.Account) (model.Account, error) {
	err := rep.db.Create(m).Error
	return *m, handleGormError(err)
}

func (rep *gormAccountRepository) Update(m *model.Account) (model.Account, error) {
	err := rep.db.Save(m).Error
	return *m, handleGormError(err)
}

func (rep *gormAccountRepository) Delete(m *model.Account) error {
	err := rep.db.Delete(m).Error
	return handleGormError(err)
}

func (rep *gormAccountRepository) InsertTx(m *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Create(m).Error
	return *m, handleGormError(err)
}

func (rep *gormAccountRepository) UpdateTx(m *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Save(m).Error
	return *m, handleGormError(err)
}

func (rep *gormAccountRepository) DeleteTx(m *model.Account, tx *gorm.DB) error {
	err := tx.Delete(m).Error
	return handleGormError(err)
}
