package pwd

import (
	"errors"
	"time"

	"github.com/elga-io/redir/api/v1/auth/jwt"
	orm "gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

// AuthStore implements database operations for account pwd authentication.
type AuthStore struct {
	db *orm.DB
}

// NewAuthStore return an AuthStore.
func NewAuthStore(db *orm.DB) *AuthStore {
	return &AuthStore{
		db: db,
	}
}

// GetAccount returns an account by ID.
func (s *AuthStore) GetAccount(id int) (Account, error) {
	a := Account{ID: id}
	result := s.db.Model(&a).Where("id = ?", id).First(&a)
	if result.RowsAffected == 0 {
		return a, errors.New("ErrNotFound")
	}
	return a, nil
}

// GetAccountByEmail returns an account by email.
func (s *AuthStore) GetAccountByEmail(e string) (Account, error) {
	a := Account{Email: e}
	result := s.db.Model(&a).Where("email = ?", a.Email).First(&a)
	if result.RowsAffected == 0 {
		return a, errors.New("ErrNotFound")
	}
	return a, nil
}

// AuthAccount returns an account by email.
func (s *AuthStore) AuthAccount(e, p string) (Account, error) {
	var err error
	var a = Account{Email: e, Password: p}

	result := s.db.Model(&a).Where("email = ?", a.Email).First(&a)
	if result.RowsAffected == 0 {
		return a, errors.New("ErrNotFound")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(p)); err != nil {
		return a, err
	}
	return a, err
}

// UpdateAccount upates account data related to pwd authentication.
func (s *AuthStore) UpdateAccount(a Account) error {
	if err := s.db.Model(a).Update("last_login", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

// GetToken returns refresh token by token identifier.
func (s *AuthStore) GetToken(t string) (jwt.Token, error) {
	ts := jwt.Token{Token: t}
	if err := s.db.Model(&ts).Where("token = ?", ts.Token).First(&ts).Error; err != nil {
		return ts, err
	}
	return ts, nil
}

// CreateOrUpdateToken creates or updates an existing refresh token.
func (s *AuthStore) CreateOrUpdateToken(t *jwt.Token) error {
	var err error
	if t.ID == 0 {
		err = s.db.Model(&t).Create(&t).Error
	} else {
		err = s.db.Model(&t).Save(&t).Error
	}
	return err
}

// DeleteToken deletes a refresh token.
func (s *AuthStore) DeleteToken(t *jwt.Token) error {
	err := s.db.Delete(t).Error
	return err
}

// PurgeExpiredToken deletes expired refresh token.
func (s *AuthStore) PurgeExpiredToken() error {
	t := jwt.Token{}
	err := s.db.Model(&t).Where("expiry < ?", time.Now()).Delete(&t).Error
	return err
}
