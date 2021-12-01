package entity

import (
	"time"
)

// Link represents a link record.
type Link struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	URLShort string `json:"url_short" gorm:"index"`
	URLFull  string `json:"url_full"`
	Title    string `json:"title"`
	Active   string `json:"active"`

	UserID string `json:"-" gorm:"index"`
}

// GetID returns the Link ID.
func (l Link) GetID() string {
	return l.ID
}
