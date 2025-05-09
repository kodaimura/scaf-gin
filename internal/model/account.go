package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Id        int            `db:"id" json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `db:"name" json:"name"`
	Password  string         `db:"password" json:"password"`
	CreatedAt time.Time      `db:"created_at" json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at" json:"deleted_at" gorm:"index"`
}
