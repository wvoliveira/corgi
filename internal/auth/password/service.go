package password

import (
	"context"
	"errors"
	"github.com/elga-io/corgi/internal/entity"
	e "github.com/elga-io/corgi/pkg/errors"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, email, password string) (entity.Identity, entity.User, error)
	Register(ctx context.Context, email, password string) error

	HTTPLogin(c *gin.Context)
	HTTPRegister(c *gin.Context)
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetUID returns the e-mail, google id, facebook id, etc.
	GetUID() string
	// GetRole returns the role.
	GetRole() string
}

type service struct {
	logger          log.Logger
	db              *gorm.DB
	signingKey      string
	tokenExpiration int
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, signingKey string, tokenExpiration int, logger log.Logger) Service {
	return service{logger, db, signingKey, tokenExpiration}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, email, password string) (identity entity.Identity, user entity.User, err error) {
	logger := s.logger.With(ctx, "email", email)

	err = s.db.Debug().Model(&user).Where("provider = ? AND uid = ?", "email", email).Association("Identities").Find(&identity)

	// to delete.
	logger.Debug("user", user)
	logger.Debug("identity", identity)

	//err = s.db.Model(&entity.Identity{}).Where("provider = ? AND uid = ?", "email", email).First(&identity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return identity, user, err
	} else if err != nil {
		return identity, user, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(identity.Password), []byte(password)); err != nil {
		logger.Infof("authentication failed")
		return identity, user, e.ErrUnauthorized
	}

	logger.Infof("authentication successful")

	accessToken, err := s.generateAccessToken(identity, user)
	if err != nil {
		return identity, user, errors.New("error to generate access token: " + err.Error())
	}

	tokenID, refreshToken, err := s.generateRefreshToken(identity, user)
	if err != nil {
		return identity, user, errors.New("error to generate refresh token: " + err.Error())
	}

	token := entity.Token{
		ID:           tokenID,
		CreatedAt:    time.Now(),
		RefreshToken: refreshToken,
		UserID:       identity.ID,
	}

	err = s.db.Debug().Model(&entity.Token{}).Create(&token).Error

	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return
}

// Register a new user to our database.
func (s service) Register(ctx context.Context, email, password string) (err error) {
	identity := entity.Identity{}
	user := entity.User{}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return e.ErrInternalServerError
	}

	identity.ID = uuid.New().String()
	identity.CreatedAt = time.Now()
	identity.Provider = "email"
	identity.UID = email
	identity.Password = string(hashedPassword)

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.Role = "user"
	user.Active = "true"
	user.Identities = append(user.Identities, identity)

	err = s.db.Debug().Model(&entity.User{}).Create(&user).Error
	return
}

func (s service) generateAccessToken(identity entity.Identity, user entity.User) (at string, err error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = identity.GetID()
	claims["uid"] = identity.GetUID()
	claims["role"] = user.GetRole()
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	at, err = accessToken.SignedString([]byte(s.signingKey))
	if err != nil {
		err = errors.New("error to generate access token: " + err.Error())
		return
	}
	return
}

func (s service) generateRefreshToken(identity entity.Identity, user entity.User) (id, rt string, err error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	claims := refreshToken.Claims.(jwt.MapClaims)
	id = uuid.New().String()

	claims["id"] = id
	claims["sub"] = 1
	claims["user_id"] = identity.GetID()
	claims["user_uid"] = identity.GetUID() // e-mail, google id, facebook id, etc
	claims["user_role"] = user.GetRole()
	claims["exp"] = time.Now().AddDate(0, 0, 7).Unix()

	rt, err = refreshToken.SignedString([]byte(s.signingKey))
	if err != nil {
		err = errors.New("error to generate refresh token: " + err.Error())
		return
	}
	return
}
