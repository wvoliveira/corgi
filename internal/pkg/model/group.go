package model

import (
	"database/sql"
	"time"
)

type Group struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`

	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`

	// CreatedBy and OnwerID:
	// UserID but you can pass the owner to another user
	CreatedBy string `json:"created_by"`
	OwnerID   string `json:"owner_id"`
}

type GroupInvite struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`

	GroupID   string `json:"group_id"`
	UserID    string `json:"user_id"`
	InvitedBy string `json:"invited_by"`
}
