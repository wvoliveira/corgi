package google

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/elga-io/corgi/internal/config"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/elga-io/corgi/pkg/jwt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
	"io/ioutil"
	"net/url"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, redirectURL string) (string, error)
	Callback(ctx context.Context, callbackURL string, r callbackRequest) (entity.Token, error)

	HTTPLogin(c *gin.Context)
	HTTPCallback(c *gin.Context)

	Routers(r *gin.Engine)
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
	logger  log.Logger
	db      *gorm.DB
	cfg     config.Config
	store   cookie.Store
	enforce *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, cfg config.Config, store cookie.Store, enforce *casbin.Enforcer) Service {
	return service{logger, db, cfg, store, enforce}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, callbackURL string) (redirectURL string, err error) {
	logger := s.logger.With(ctx)
	logger.Info("creating a URL for send back to redirect user")

	conf := &oauth2.Config{
		ClientID:     s.cfg.Auth.Google.ClientID,
		ClientSecret: s.cfg.Auth.Google.ClientSecret,
		RedirectURL:  callbackURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	redirectURL = conf.AuthCodeURL("state")
	return
}

func (s service) Callback(ctx context.Context, callbackURL string, r callbackRequest) (token entity.Token, err error) {
	logger := s.logger.With(ctx)
	logger.Info("Callback func to get token from Google")

	conf := &oauth2.Config{
		ClientID:     s.cfg.Auth.Google.ClientID,
		ClientSecret: s.cfg.Auth.Google.ClientSecret,
		RedirectURL:  callbackURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	oauthToken, err := conf.Exchange(ctx, r.Code)
	if err != nil {
		logger.Error("error to exchange token from Google", err.Error())
		return
	}

	client := conf.Client(ctx, oauthToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(oauthToken.AccessToken))
	if err != nil {
		logger.Error("error to get userinfo email", err.Error())
		return
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error to read body response", err.Error())
		return
	}

	googleUser := entity.GoogleUserInfo{}
	err = json.Unmarshal(response, &googleUser)
	if err != nil {
		logger.Error("error to unmarshal userinfo from google", err.Error())
		return
	}

	identity := entity.Identity{}
	// check if user exists in database.
	err = s.db.Debug().Model(entity.Identity{}).Where("provider = ? AND UID = ?", "google", googleUser.ID).First(&identity).Error
	if err == gorm.ErrRecordNotFound {
		identity := entity.Identity{}
		user := entity.User{}

		identity.ID = uuid.New().String()
		identity.CreatedAt = time.Now()
		identity.LastLogin = identity.CreatedAt
		identity.Provider = "google"
		identity.UID = googleUser.ID
		identity.Verified = &googleUser.VerifiedEmail
		identity.VerifiedAt = identity.CreatedAt

		t := true
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.Name = googleUser.Name
		user.Role = "user"
		user.Active = &t
		user.Identities = append(user.Identities, identity)

		err = s.db.Debug().Model(&entity.User{}).Create(&user).Error
		if err != nil {
			logger.Error("error to create google user", err.Error())
			return
		}
	} else if err != nil {
		logger.Error("error to get user from database", err.Error())
	}

	// Get user info.
	if identity.UserID == "" {
		err = s.db.Debug().Model(&entity.Identity{}).Where("provider = ? AND uid = ?", "google", googleUser.ID).First(&identity).Error
		if err != nil {
			logger.Error("error to get identity from database", err.Error())
			return
		}
	}

	user := entity.User{}
	err = s.db.Debug().Model(&entity.User{}).Where("id = ?", identity.UserID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return token, err
	} else if err != nil {
		return token, err
	}

	accessToken, err := jwt.GenerateAccessToken(s.cfg.App.SecretKey, identity, user)
	if err != nil {
		return token, errors.New("error to generate access token: " + err.Error())
	}

	refreshToken, err := jwt.GenerateRefreshToken(s.cfg.App.SecretKey, identity, user)
	if err != nil {
		return token, errors.New("error to generate refresh token: " + err.Error())
	}

	refreshToken.UserID = identity.UserID
	err = s.db.Debug().Model(&entity.Token{}).Create(&refreshToken).Error
	if err != nil {
		return
	}

	token.ID = refreshToken.ID
	token.AccessToken = accessToken.AccessToken
	token.RefreshToken = refreshToken.RefreshToken
	token.AccessExpires = accessToken.AccessExpires
	return
}
