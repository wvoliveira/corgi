package entity

import "time"

// Token represent a JWT struct.
type Token struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastUse   time.Time `json:"last_use"`

	Token     string    `json:"token"`
	ExpiresIn time.Time `json:"expires_in"`

	UserID string `json:"user_id" gorm:"index"`
}
