package model

import "time"

type PasswordResetToken struct {
	ID        int        `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AccountID int        `db:"account_id" gorm:"column:account_id"`
	TokenHash string     `db:"token_hash" gorm:"column:token_hash;uniqueIndex"`
	ExpiresAt time.Time  `db:"expires_at" gorm:"column:expires_at"`
	UsedAt    *time.Time `db:"used_at" gorm:"column:used_at"`
	CreatedAt time.Time  `db:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time  `db:"updated_at" gorm:"column:updated_at"`
}
