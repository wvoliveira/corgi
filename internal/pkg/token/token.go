package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

const (
	tokenAuth     = "token_auth:%s"
	tokenPersonal = "token_personal:%s"
)

var jwtKey = []byte(viper.GetString("SECRET_KEY"))

type JWTClaim struct {
	UserID string
	model.User
	jwt.StandardClaims
}

func CreateToken(c context.Context, cache *redis.Client, typeToken string, userID string, rememberMe bool) (token string, err error) {
	if typeToken != "auth" && typeToken != "personal" {
		return token, e.ErrTokenType
	}

	tokens := map[string]string{
		"auth":     tokenAuth,
		"personal": tokenPersonal,
	}

	expirationTime := 48 * time.Hour
	if typeToken == "personal" || rememberMe {
		expirationTime = 0
	}

	token = ulid.Make().String()

	key := fmt.Sprintf(tokens[typeToken], token)
	status := cache.Set(c, key, userID, expirationTime)

	return token, status.Err()
}

func GenerateJWTAccess(user model.User) (accessToken, refreshToken string, err error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &JWTClaim{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString(jwtKey)
	if err != nil {
		return
	}

	expirationTime = time.Now().Add(8640 * time.Hour)

	claims.User = model.User{}
	claims.UserID = user.ID
	claims.StandardClaims.ExpiresAt = expirationTime.Unix()

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err = token.SignedString(jwtKey)

	return
}

func GenerateJWTRefresh(user model.User) (tokenString string, err error) {
	expirationTime := time.Now().Add(8640 * time.Hour)

	claims := &JWTClaim{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (user model.User, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token is expired")
		return
	}

	user = claims.User
	return
}
