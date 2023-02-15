package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Identity struct {
	ID        string    `json:"id" gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login"`

	Provider string `json:"provider"` // phone, email, wechat, github...
	UID      string `json:"uid"`      // username, e-mail, google id, facebook id, etc
	Password string `json:"password"`
	UserID   string `json:"user_id"`

	Verified   *bool     `json:"verified" gorm:"default:false"`
	VerifiedAt time.Time `json:"confirmed_at"`
}

// IdentityInfo pass these info in middleware.
type IdentityInfo struct {
	ID             string
	Provider       string
	UID            string
	UserID         string
	UserRole       string
	RefreshTokenID string
}

func (iden *Identity) CheckPassword(password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(iden.Password))
}

func (iden *Identity) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	// Store hashed password in own identity.
	iden.Password = string(bytes)
	return nil
}
