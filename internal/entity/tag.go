package entity

import "time"

type Tag struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Users []User `gorm:"many2many:user_tags;"`
}
