package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID           int64          `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Email        *string        `db:"email" gorm:"column:email;type:text;uniqueIndex"`
	LoginID      string         `db:"login_id" gorm:"column:login_id;type:text;uniqueIndex"`
	PasswordHash string         `db:"password_hash" gorm:"column:password_hash;type:text"`
	TokenVersion int            `db:"token_version" gorm:"column:token_version;default:1"`
	FirstName    string         `db:"first_name" gorm:"column:first_name;size:100"`
	LastName     string         `db:"last_name" gorm:"column:last_name;size:100"`
	DisabledAt   *time.Time     `db:"disabled_at" gorm:"column:disabled_at;type:timestamptz"`
	DeletedAt    gorm.DeletedAt `db:"deleted_at" gorm:"column:deleted_at;type:timestamptz"`
	CreatedAt    time.Time      `db:"created_at" gorm:"column:created_at;type:timestamptz"`
	UpdatedAt    time.Time      `db:"updated_at" gorm:"column:updated_at;type:timestamptz"`
}
