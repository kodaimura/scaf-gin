package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID           int            `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Email        *string        `db:"email" gorm:"column:email;uniqueIndex"`
	LoginID      string         `db:"login_id" gorm:"column:login_id;uniqueIndex"`
	PasswordHash string         `db:"password_hash" gorm:"column:password_hash"`
	TokenVersion int            `db:"token_version" gorm:"column:token_version;default:1"`
	FirstName    string         `db:"first_name" gorm:"column:first_name"`
	LastName     string         `db:"last_name" gorm:"column:last_name"`
	DisabledAt   *time.Time     `db:"disabled_at" gorm:"column:disabled_at"`
	CreatedAt    time.Time      `db:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time      `db:"updated_at" gorm:"column:updated_at"`
	DeletedAt    gorm.DeletedAt `db:"deleted_at" gorm:"column:deleted_at;index"`
}
