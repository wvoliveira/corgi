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
