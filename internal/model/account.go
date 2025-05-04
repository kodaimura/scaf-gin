package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	AccountId       int            `db:"account_id" json:"account_id" gorm:"primaryKey;autoIncrement"`
	AccountName     string         `db:"account_name" json:"account_name"`
	AccountPassword string         `db:"account_password" json:"account_password"`
	CreatedAt       time.Time      `db:"created_at" json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time      `db:"updated_at" json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `db:"deleted_at" json:"deleted_at" gorm:"index"`
}