package server

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AddAccount create a new account.
func (s Service) AddAccount(auth, payload Account) (account Account, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	if err != nil {
		return account, ErrInternalServerError
	}

	payload.ID = uuid.New().String()
	payload.Password = string(hashedPassword)
	payload.CreatedAt = time.Now()
	payload.Active = "true"

	account = payload
	err = s.db.Model(&Account{}).Create(&account).Error
	return
}

// FindAccountByID find a account with specific ID.
func (s Service) FindAccountByID(_ Account, id string) (account Account, err error) {
	err = s.db.Model(&Account{}).Where("id = ?", id).First(&account).Error
	return
}

// FindAccounts Get a list of accounts.
func (s Service) FindAccounts(_ Account, offset, limit int) (accounts []Account, err error) {
	err = s.db.Model(&Account{}).Limit(limit).Offset(offset).Find(&accounts).Error
	return
}

// UpdateAccount update specific account fields.
func (s Service) UpdateAccount(_ Account, id string, payload Account) (err error) {
	err = s.db.Model(&Account{}).Where("id = ?", id).Updates(&payload).Error
	return
}

// DeleteAccount delete specific account by ID.
func (s Service) DeleteAccount(_ Account, id string) (err error) {
	var account Account
	err = s.db.Model(&Account{}).Where("id = ?", id).Delete(&account).Error
	return
}
