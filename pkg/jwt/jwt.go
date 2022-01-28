package jwt

import (
	"errors"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

// ValidToken verify if token is a valid one.
func ValidToken(secret string, payload string) (claims jwt.MapClaims, valid bool) {
	var SigningKey = []byte(secret)

	token, err := jwt.Parse(payload, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return token, e.ErrParseToken
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

// UpdateAccessToken generate access and refresh tokens.
func UpdateAccessToken(secret string, claims jwt.MapClaims) (token entity.Token, err error) {
	identity := entity.Identity{
		ID:       claims["identity_id"].(string),
		Provider: claims["identity_provider"].(string),
		UID:      claims["identity_uid"].(string),
	}

	user := entity.User{
		ID:   claims["user_id"].(string),
		Role: claims["user_role"].(string),
	}

	token, err = GenerateAccessToken(secret, identity, user)
	if err != nil {
		return
	}
	return
}

// UpdateRefreshToken generate refresh token.
func UpdateRefreshToken(secret string, claims jwt.MapClaims) (token entity.Token, err error) {
	identity := entity.Identity{
		ID:       claims["identity_id"].(string),
		Provider: claims["identity_provider"].(string),
		UID:      claims["identity_uid"].(string),
	}

	user := entity.User{
		ID:   claims["user_id"].(string),
		Role: claims["user_role"].(string),
	}

	token, err = GenerateRefreshToken(secret, identity, user)
	if err != nil {
		return
	}
	return
}

// GenerateAccessToken generate a new JWT token with user info in claims.
func GenerateAccessToken(secret string, identity entity.Identity, user entity.User) (token entity.Token, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)

	// This system not has a security problem. So, the token expires in 2 hours.
	tokenExpires := time.Now().Add(time.Hour * 2)

	claims := accessToken.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["identity_id"] = identity.ID
	claims["identity_provider"] = identity.Provider // e-mail, google, facebook, etc.
	claims["identity_uid"] = identity.UID           // e-mail address, google id, facebook id, etc.
	claims["user_id"] = user.ID
	claims["user_role"] = user.Role
	claims["exp"] = tokenExpires.Unix()

	at, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		err = errors.New("error to generate access token: " + err.Error())
		return
	}
	token.CreatedAt = time.Now()
	token.Token = at
	token.ExpiresIn = tokenExpires
	token.UserID = user.ID
	return
}

// GenerateRefreshToken generate a new JWT refresh token and add user info in claims.
func GenerateRefreshToken(secret string, identity entity.Identity, user entity.User) (token entity.Token, err error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	id := uuid.New().String()

	// Refresh token expires in 7 days. But I think to increase this value.
	tokenExpires := time.Now().AddDate(0, 0, 7)

	claims := refreshToken.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["sub"] = 1
	claims["identity_id"] = identity.ID
	claims["identity_provider"] = identity.Provider // e-mail, google, facebook, etc.
	claims["identity_uid"] = identity.UID           // e-mail address, google id, facebook id, etc.
	claims["user_id"] = user.ID
	claims["user_role"] = user.Role
	claims["exp"] = tokenExpires.Unix()

	rt, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		err = errors.New("error to generate refresh token: " + err.Error())
		return
	}

	token.ID = id
	token.CreatedAt = time.Now()
	token.Token = rt
	token.ExpiresIn = tokenExpires
	token.UserID = user.ID
	return
}

// ValidateToken check if token is valid and return claims if true.
func ValidateToken(tokenHash, secret string) (claims jwt.MapClaims, err error) {
	token, err := jwt.Parse(tokenHash, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return token, e.ErrParseToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		return claims, errors.New("error to parse access token " + err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, err
	}
	return claims, errors.New("invalid token")
}
