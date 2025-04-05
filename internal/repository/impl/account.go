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

func (repo *gormAccountRepository) Get(a *model.Account) ([]model.Account, error) {
	var accounts []model.Account
	err := repo.db.Find(&accounts, a).Error
	return accounts, err
}

func (repo *gormAccountRepository) GetOne(a *model.Account) (model.Account, error) {
	var account model.Account
	err := repo.db.First(&account, a).Error
	return account, err
}

func (repo *gormAccountRepository) Insert(a *model.Account) (model.Account, error) {
	err := repo.db.Create(a).Error
	return a, err
}

func (repo *gormAccountRepository) Update(a *model.Account) (model.Account, error) {
	rerr := repo.db.Save(a).Error
	return a, err
}

func (repo *gormAccountRepository) Delete(a *model.Account) error {
	return repo.db.Delete(a).Error
}

func (repo *gormAccountRepository) InsertTx(a *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Create(a).Error
	return a, err
}

func (repo *gormAccountRepository) UpdateTx(a *model.Account, tx *gorm.DB) (model.Account, error) {
	err := tx.Save(a).Error
	return a, err
}

func (repo *gormAccountRepository) DeleteTx(a *model.Account, tx *gorm.DB) error {
	return tx.Delete(a).Error
}