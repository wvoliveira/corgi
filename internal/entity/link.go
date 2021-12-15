package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Link represents a link record.
type Link struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Domain  string `json:"domain" gorm:"index"`
	Keyword string `json:"keyword" gorm:"index"`
	URL     string `json:"url" gorm:"index"`
	Title   string `json:"title"`
	Active  string `json:"active"`

	UserID    string `json:"-" gorm:"index"`
	LinkLogID string `json:"-" gorm:"index"`
}

type LinkLog struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:null"`

	RemoteAddress         string `json:"remote_address"`
	UserAgentRaw          string `json:"user_agent_raw"`
	UserAgentFamily       string `json:"user_agent_family"`
	UserAgentOSFamily     string `json:"user_agent_os_family"`
	UserAgentDeviceFamily string `json:"user_agent_device_family"`
	Referer               string `json:"referer"`

	LinkID         string `json:"-" gorm:"index"`
	LocationIPv4ID string `json:"-" gorm:"index"`
}

func (l *LinkLog) BeforeCreate(tx *gorm.DB) (err error) {
	l.ID = uuid.New().String()
	l.CreatedAt = time.Now()
	return
}

// GetID returns the Link ID.
func (l Link) GetID() string {
	return l.ID
}
