package account

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	Id        int            `db:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `db:"name" gorm:"column:name"`
	Password  string         `db:"password" gorm:"column:password"`
	CreatedAt time.Time      `db:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time      `db:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at" gorm:"column:deleted_at;index"`
}
