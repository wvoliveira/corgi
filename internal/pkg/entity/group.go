package entity

import "time"

type Group struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name        string `json:"name"`
	Description string `json:"description"`

	UserID string `json:"user_id"`
	Users  []User `gorm:"many2many:user_tags;"`
}
