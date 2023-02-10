package model

import (
	"database/sql"
	"time"
)

// User represents a user info.
type User struct {
	ID        string       `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`

	Name   string `json:"name"`
	Role   string `json:"role"`
	Active *bool  `json:"active" gorm:"default:true"`

	Identities []Identity `json:"identities,omitempty"`
	Tokens     []Token    `json:"tokens,omitempty"`
	Links      []Link     `json:"links,omitempty"`

	Tags []Tag `gorm:"many2many:user_tags;" json:"tags,omitempty"`
}

type UserGoogle struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

type UserFacebook struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
