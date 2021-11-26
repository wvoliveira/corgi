package entity

import (
	"time"
)

// Link represents a link record.
type Link struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Keyword     string `json:"keyword" gorm:"index"`
	Destination string `json:"destination"`
	Title       string `json:"title"`
	Active      string `json:"active"`

	UserID string `json:"user_id" gorm:"index"`
}

// GetID returns the Link ID.
func (l Link) GetID() string {
	return l.ID
}
