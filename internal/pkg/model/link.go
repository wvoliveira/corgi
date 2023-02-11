package model

import (
	"database/sql"
	"time"
)

// Link represents a link record.
type Link struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`

	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Active  string `json:"active"`

	UserID string `json:"-"`
}

// LinkLog model to store redirects logs.
type LinkLog struct {
	ID        string    `json:"id"`
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

	LinkID string `json:"link_id"`
}
