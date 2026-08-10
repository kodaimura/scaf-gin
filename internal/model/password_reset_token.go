package model

import "time"

type PasswordResetToken struct {
	ID        int64      `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	AccountID int64      `db:"account_id" gorm:"column:account_id;index"`
	TokenHash string     `db:"token_hash" gorm:"column:token_hash;type:text;uniqueIndex"`
	ExpiresAt time.Time  `db:"expires_at" gorm:"column:expires_at;type:timestamptz"`
	UsedAt    *time.Time `db:"used_at" gorm:"column:used_at;type:timestamptz"`
	CreatedAt time.Time  `db:"created_at" gorm:"column:created_at;type:timestamptz"`
	UpdatedAt time.Time  `db:"updated_at" gorm:"column:updated_at;type:timestamptz"`
}
