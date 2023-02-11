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
	Users     []User `json:"users,omitempty"`
}
