package server

import (
	"gorm.io/gorm"
	"time"
)

// AuthIdentity auth identity session model
type AuthIdentity struct {
	gorm.Model
	Provider    string // phone, email, wechat, github...
	UID         string `gorm:"column:uid"`
	Password    string
	UserID      string
	ConfirmedAt *time.Time
}
