package server

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login login with email and password.
func (s Service) Login(payload Account) (account Account, err error) {
	err = s.db.Model(&Account{}).Where("email = ?", payload.Email).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return account, err
	} else if err != nil {
		return account, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password)); err != nil {
		return account, ErrUnauthorized
	}

	accessToken, err := s.generateAccessToken(account)
	if err != nil {
		return account, errors.New("error to generate access token: " + err.Error())
	}

	tokenID, refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return account, errors.New("error to generate refresh token: " + err.Error())
	}

	token := Token{
		ID:           tokenID,
		CreatedAt:    time.Now(),
		RefreshToken: refreshToken,
		AccountID:    account.ID,
	}

	err = s.db.Debug().Model(&Token{}).Create(&token).Error

	account.AccessToken = accessToken
	account.RefreshToken = refreshToken
	return
}

// Register register with e-mail and password.
func (s Service) Register(payload Account) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	if err != nil {
		return ErrInternalServerError
	}

	payload.ID = uuid.New().String()
	payload.CreatedAt = time.Now()
	payload.Password = string(hashedPassword)
	payload.Role = "user"
	payload.Active = "true"
	err = s.db.Debug().Model(&Account{}).Create(&payload).Error
	return
}
