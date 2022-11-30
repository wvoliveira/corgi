package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Link represents a link record.
type Link struct {
	ID        string     `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`

	Domain  string `json:"domain" gorm:"index"`
	Keyword string `json:"keyword" gorm:"index"`
	URL     string `json:"url" gorm:"index"`
	Title   string `json:"title"`
	Active  string `json:"active"`

	UserID string `json:"-" gorm:"index"`
}

// LinkLog model to store redirections logs.
type LinkLog struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`

	// Log object content.
	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	URL     string `json:"url"`
	Title   string `json:"title"`

	// Client "requester" content.
	RemoteAddress         string `json:"remote_address"`
	UserAgent             string `json:"user_agent"`
	UserAgentFamily       string `json:"user_agent_family"`
	UserAgentOSFamily     string `json:"user_agent_os_family"`
	UserAgentDeviceFamily string `json:"user_agent_device_family"`
	Referer               string `json:"referer"`

	LinkID string `json:"link_id" gorm:"index"`
}

func (l *Link) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	l.CreatedAt = time.Now()
	return
}

func (l *Link) BeforeUpdate(tx *gorm.DB) (err error) {
	t := time.Now()
	l.UpdatedAt = &t
	return
}

func (l *LinkLog) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	return
}

// GetID returns the Link ID.
func (l Link) GetID() string {
	return l.ID
}
