package pwd

import (
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"gorm.io/gorm"

	"github.com/elga-io/redir/api/v1/auth/jwt"
)

// Account represents an authenticated application user
type Account struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	LastLogin time.Time `json:"last_login,omitempty"`

	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"-"`
	Token    []jwt.Token `json:"token,omitempty"`
	Active   *bool       `json:"active" gorm:"type:bool;default:true" example:"false"`
	Roles    []string    `json:"roles" gorm:"array"`

	UserID string `json:"user_id"`
}

// BeforeInsert hook executed before database insert operation.
func (a *Account) BeforeInsert(db gorm.DB) error {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
		a.UpdatedAt = now
	}
	return a.Validate()
}

// BeforeUpdate hook executed before database update operation.
func (a *Account) BeforeUpdate(db gorm.DB) error {
	a.UpdatedAt = time.Now()
	return a.Validate()
}

// BeforeDelete hook executed before database delete operation.
func (a *Account) BeforeDelete(db gorm.DB) error {
	return nil
}

// Validate validates Account struct and returns validation errors.
func (a *Account) Validate() error {
	a.Email = strings.TrimSpace(a.Email)
	a.Email = strings.ToLower(a.Email)

	return validation.ValidateStruct(a,
		validation.Field(&a.Email, validation.Required, is.Email, is.LowerCase),
	)
}

// CanLogin returns true if user is allowed to login.
func (a *Account) CanLogin() *bool {
	return a.Active
}

// Claims returns the account's claims to be signed
func (a *Account) Claims() jwt.AppClaims {
	return jwt.AppClaims{
		ID:    a.ID,
		Sub:   a.Name,
		Roles: a.Roles,
	}
}

// CheckPassword returns true if user is allowed to login.
func (a *Account) CheckPassword() *bool {
	return a.Active
}
