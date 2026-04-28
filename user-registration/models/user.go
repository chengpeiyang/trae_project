package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(32);not null" json:"-"`
	Email     string         `gorm:"type:varchar(100)" json:"email"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	Status    int            `gorm:"type:tinyint;default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
