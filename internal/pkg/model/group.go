package model

import (
	"database/sql"
	"time"
)

type Group struct {
	ID            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
	UpdatedAtNull sql.NullTime `json:"-"`

	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`

	// CreatedBy and OwnerID:
	// UserID but you can pass the owner to another user
	CreatedBy string `json:"created_by"`
	OwnerID   string `json:"owner_id"`
}

type GroupInvite struct {
	ID            string       `json:"id"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at"`
	UpdatedAtNull sql.NullTime `json:"-"`

	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	InvitedBy string `json:"invited_by"`
	Accepted  bool   `json:"accepted"`
}
