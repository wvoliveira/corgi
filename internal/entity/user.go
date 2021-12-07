package entity

import "time"

// User represents a user info.
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Active *bool  `json:"active" gorm:"default:true"`

	Identities []Identity
	Tokens     []Token
	Links      []Link

	Tags []Tag `gorm:"many2many:user_tags;"`
}

// GetID returns the user ID.
func (u User) GetID() string {
	return u.ID
}

// GetRole returns the role.
func (u User) GetRole() string {
	return u.Role
}
