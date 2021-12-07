package facebook

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/elga-io/corgi/internal/config"
	"github.com/elga-io/corgi/internal/entity"
	"github.com/elga-io/corgi/pkg/jwt"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
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
	logger log.Logger
	db     *gorm.DB
	cfg    config.Config
	store  cookie.Store
}

// NewService creates a new authentication service.
func NewService(logger log.Logger, db *gorm.DB, cfg config.Config, store cookie.Store) Service {
	return service{logger, db, cfg, store}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, callbackURL string) (redirectURL string, err error) {
	logger := s.logger.With(ctx)
	logger.Info("creating a URL for send back to redirect user")

	conf := &oauth2.Config{
		ClientID:     s.cfg.Auth.Facebook.ClientID,
		ClientSecret: s.cfg.Auth.Facebook.ClientSecret,
		RedirectURL:  callbackURL,
		Scopes: []string{
			"public_profile",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
	// Redirect user to Facebook's consent page to ask for permission
	// for the scopes specified above.
	redirectURL = conf.AuthCodeURL("state")
	return
}

func (s service) Callback(ctx context.Context, callbackURL string, r callbackRequest) (token entity.Token, err error) {
	logger := s.logger.With(ctx)
	logger.Info("Callback func to get token from Facebook")

	conf := &oauth2.Config{
		ClientID:     s.cfg.Auth.Facebook.ClientID,
		ClientSecret: s.cfg.Auth.Facebook.ClientSecret,
		RedirectURL:  callbackURL,
		Scopes: []string{
			"public_profile",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}

	oauthToken, err := conf.Exchange(ctx, r.Code)
	if err != nil {
		logger.Error("error to exchange token from Facebook", err.Error())
		return
	}

	client := conf.Client(ctx, oauthToken)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + url.QueryEscape(oauthToken.AccessToken))
	if err != nil {
		logger.Error("error to get userinfo email", err.Error())
		return
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error to read body response", err.Error())
		return
	}

	type facebookUserInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	facebookUser := facebookUserInfo{}
	err = json.Unmarshal(response, &facebookUser)
	if err != nil {
		logger.Error("error to unmarshal userinfo from Facebook", err.Error())
		return
	}

	identity := entity.Identity{}
	// check if user exists in database.
	err = s.db.Debug().Model(entity.Identity{}).Where("provider = ? AND UID = ?", "facebook", facebookUser.ID).First(&identity).Error
	if err == gorm.ErrRecordNotFound {
		identity := entity.Identity{}
		user := entity.User{}

		identity.ID = uuid.New().String()
		identity.CreatedAt = time.Now()
		identity.LastLogin = identity.CreatedAt
		identity.Provider = "facebook"
		identity.UID = facebookUser.ID
		//identity.Verified = false
		//identity.VerifiedAt = identity.CreatedAt

		t := true
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.Name = facebookUser.Name
		user.Role = "user"
		user.Active = &t
		user.Identities = append(user.Identities, identity)

		err = s.db.Debug().Model(&entity.User{}).Create(&user).Error
		if err != nil {
			logger.Error("error to create Facebook user", err.Error())
			return
		}
	} else if err != nil {
		logger.Error("error to get user from database", err.Error())
	}

	// Get user info.
	if identity.UserID == "" {
		err = s.db.Debug().Model(&entity.Identity{}).Where("provider = ? AND uid = ?", "facebook", facebookUser.ID).First(&identity).Error
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
