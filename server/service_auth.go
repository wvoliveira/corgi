package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SignIn login with email and password.
func (s Service) SignIn(payload Account) (account Account, err error) {
	err = s.db.Model(&Account{}).Where("email = ?", payload.Email).First(&account).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return account, err
	} else if err != nil {
		return account, err
	}

	fmt.Println("PAYLOAD")
	fmt.Println(payload)

	fmt.Println("ACCOUNT")
	fmt.Println(account)

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password)); err != nil {
		return account, ErrUnauthorized
	}

	tokenHash, err := generateJWT(s.secret, account)
	if err != nil {
		return
	}

	account.Token = tokenHash
	return
}

// SignUp register with e-mail and password.
func (s Service) SignUp(payload Account) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 8)
	if err != nil {
		return ErrInternalServerError
	}

	payload.ID = uuid.New().String()
	payload.CreatedAt = time.Now()
	payload.Password = string(hashedPassword)
	payload.Role = "user"
	payload.Active = "true"
	err = s.db.Model(&Account{}).Create(&payload).Error
	return
}

func generateJWT(secretKey string, a Account) (hash string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = a.ID
	claims["email"] = a.Email
	claims["role"] = a.Role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	hash, err = token.SignedString([]byte(secretKey))
	if err != nil {
		err = errors.New("error to generate JWT: " + err.Error())
		return
	}

	return
}
