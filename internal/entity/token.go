package entity

import "time"

// Token represent a JWT struct.
type Token struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastUse   time.Time `json:"last_use"`

	AccessToken  string `json:"access_token" gorm:"-"`
	RefreshToken string `json:"refresh_token"`

	AccessExpires int64 `json:"access_expires" gorm:"-"`
	RefreshExpires int64 `json:"refresh_expires"`

	UserID string `json:"user_id" gorm:"index"`
}

// GetID returns the Token ID.
func (t Token) GetID() string {
	return t.ID
}
