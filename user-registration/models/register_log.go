package models

import (
	"time"
)

type RegisterLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	Username      string    `gorm:"type:varchar(50);not null" json:"username"`
	IP            string    `gorm:"type:varchar(50);not null" json:"ip"`
	UserAgent     string    `gorm:"type:varchar(500)" json:"user_agent"`
	Source        string    `gorm:"type:varchar(50)" json:"source"`
	DeviceType    string    `gorm:"type:varchar(50)" json:"device_type"`
	OS            string    `gorm:"type:varchar(50)" json:"os"`
	Browser       string    `gorm:"type:varchar(50)" json:"browser"`
	Referrer      string    `gorm:"type:varchar(255)" json:"referrer"`
	RequestURI    string    `gorm:"type:varchar(255)" json:"request_uri"`
	CreatedAt     time.Time `json:"created_at"`
}

func (RegisterLog) TableName() string {
	return "register_logs"
}
