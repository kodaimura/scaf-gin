package repository

import (
	"scaf-gin/internal/model"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Get(m *model.Account) ([]model.Account, error)
	GetOne(m *model.Account) (model.Account, error)
	GetAll(m *model.Account) ([]model.Account, error)

	Insert(m *model.Account) (model.Account, error)
	Update(m *model.Account) (model.Account, error)
	Delete(m *model.Account) error

	InsertTx(m *model.Account, tx *gorm.DB) (model.Account, error)
	UpdateTx(m *model.Account, tx *gorm.DB) (model.Account, error)
	DeleteTx(m *model.Account, tx *gorm.DB) error
}
