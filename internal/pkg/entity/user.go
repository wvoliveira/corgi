package entity

import "time"

// User represents a user info.
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   string `json:"name"`
	Role   string `json:"role"`
	Active *bool  `json:"active" gorm:"default:true"`

	Identities []Identity
	Tokens     []Token
	Links      []Link

	Groups []Group `gorm:"many2many:user_groups;"`
	Tags   []Tag   `gorm:"many2many:user_tags;"`
}
