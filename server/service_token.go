package server

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Refresh refresh access token.
func (s Service) Refresh(a Account, payload Token) (token Token, err error) {
	claims, ok := s.validToken(payload.RefreshToken)
	if !ok {
		return token, errors.New("invalid refresh token: " + err.Error())
	}

	payload.ID = claims["id"].(string)

	if err = s.db.Debug().Model(&Token{}).Where("id = ?", payload.ID).First(&token).Error; err != nil {
		return token, errors.New("error to get refresh token from database: " + err.Error())
	}

	_, ok = s.validToken(token.RefreshToken)

	if !ok {
		return token, errors.New("invalid token from database: " + err.Error())
	}

	at, err := s.generateAccessToken(a)
	if err != nil {
		return token, errors.New("error to generate access token: " + err.Error())
	}
	token.AccessToken = at
	return
}

func (s Service) validToken(payload string) (claims jwt.MapClaims, valid bool) {
	var SigningKey = []byte(s.secret)

	token, err := jwt.Parse(payload, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return token, ErrParseToken
		}
		return SigningKey, nil
	})

	if err != nil {
		return claims, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	}
	return
}

func (s Service) generateAccessToken(a Account) (at string, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = a.ID
	claims["email"] = a.Email
	claims["role"] = a.Role
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	at, err = accessToken.SignedString([]byte(s.secret))
	if err != nil {
		err = errors.New("error to generate access token: " + err.Error())
		return
	}
	return
}

func (s Service) generateRefreshToken() (id, rt string, err error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims := refreshToken.Claims.(jwt.MapClaims)
	id = uuid.New().String()

	claims["id"] = id
	claims["sub"] = 1
	claims["exp"] = time.Now().AddDate(0, 0, 7).Unix()

	rt, err = refreshToken.SignedString([]byte(s.secret))
	if err != nil {
		err = errors.New("error to generate refresh token: " + err.Error())
		return
	}
	return
}
