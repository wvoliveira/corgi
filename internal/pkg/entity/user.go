package entity

import "time"

// User represents a user info.
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`

	Name   string `json:"name"`
	Role   string `json:"-"`
	Active *bool  `json:"-" gorm:"default:true"`

	Identities []Identity `json:"-"`
	Tokens     []Token    `json:"-"`
	Links      []Link     `json:"-"`

	Groups []Group `json:"-" gorm:"many2many:user_groups;"`
	Tags   []Tag   `json:"-" gorm:"many2many:user_tags;"`
}
