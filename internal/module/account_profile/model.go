package account_profile

import (
	"time"
)

type AccountProfile struct {
	AccountId   int        `db:"account_id" gorm:"column:account_id;primaryKey;autoIncrement:false"`
	DisplayName string     `db:"display_name" gorm:"column:display_name"`
	Bio         string     `db:"bio" gorm:"column:bio"`
	AvatarURL   string     `db:"avatar_url" gorm:"column:avatar_url"`
	CreatedAt   time.Time  `db:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time  `db:"updated_at" gorm:"column:updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" gorm:"column:deleted_at;index"`
}
