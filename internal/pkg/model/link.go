package model

import (
	"database/sql"
	"time"
)

// Link represents a link record.
type Link struct {
	ID            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
	UpdatedAtNull sql.NullTime `json:"-"`

	Domain  string `json:"domain"`
	Keyword string `json:"keyword"`
	URL     string `json:"url"`
	Title   string `json:"title"`
	Active  string `json:"active"`

	UserID string `json:"-"`
}
