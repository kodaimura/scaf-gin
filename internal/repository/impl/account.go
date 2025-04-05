package impl

import (
	"goscaf/internal/model"
	"gorm.io/gorm"
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
	return accounts, err
}

func (rep *gormAccountRepository) GetOne(a *model.Account) (model.Account, error) {
	var account model.Account
	err := rep.db.First(&account, a).Error
	return account, err
}

func (rep *gormAccountRepository) Insert(a *model.Account) (model.Account, error) {
	err := rep.db.Create(a).Error
	return a, err
}

func (rep *gormAccountRepository) Update(a *model.Account) (model.Account, error) {
	rerr := rep.db.Save(a).Error
	return a, err
}

func (rep *gormAccountRepository) Delete(a *model.Account) error {
	return rep.db.Delete(a).Error
}

func (rep *gormAccountRepository) InsertTx(a *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Create(a).Error
	return a, err
}

func (rep *gormAccountRepository) UpdateTx(a *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Save(a).Error
	return a, err
}

func (rep *gormAccountRepository) DeleteTx(a *model.Account, tx *gorm.DB) error {
	return tx.Delete(a).Error
}