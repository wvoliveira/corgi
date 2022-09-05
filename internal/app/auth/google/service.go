package google

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/wvoliveira/corgi/internal/app/config"
	"github.com/wvoliveira/corgi/internal/app/entity"
	"github.com/wvoliveira/corgi/internal/pkg/jwt"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(ctx context.Context, redirectURL string) (string, error)
	Callback(ctx context.Context, callbackURL string, r callbackRequest) (entity.Token, entity.Token, error)

	NewHTTP(r *mux.Router)
	HTTPLogin(w http.ResponseWriter, r *http.Request)
	HTTPCallback(w http.ResponseWriter, r *http.Request)
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
	db      *gorm.DB
	cfg     config.Config
	store   *sessions.CookieStore
	enforce *casbin.Enforcer
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, cfg config.Config, store *sessions.CookieStore, enforce *casbin.Enforcer) Service {
	return service{db, cfg, store, enforce}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, callbackURL string) (redirectURL string, err error) {
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

func (s service) Callback(ctx context.Context, callbackURL string, r callbackRequest) (tokenAccess, tokenRefresh entity.Token, err error) {
	l := logger.Logger(ctx)

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
		l.Error().Caller().Msg(err.Error())
		return
	}

	client := conf.Client(ctx, oauthToken)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(oauthToken.AccessToken))
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	googleUser := entity.GoogleUserInfo{}
	err = json.Unmarshal(response, &googleUser)
	if err != nil {
		l.Error().Caller().Msg(err.Error())
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
			l.Error().Caller().Msg(err.Error())
			return
		}
	} else if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	// Get user info.
	if identity.UserID == "" {
		err = s.db.Debug().Model(&entity.Identity{}).Where("provider = ? AND uid = ?", "google", googleUser.ID).First(&identity).Error
		if err != nil {
			l.Error().Caller().Msg(err.Error())
			return
		}
	}

	user := entity.User{}
	err = s.db.Debug().Model(&entity.User{}).Where("id = ?", identity.UserID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tokenAccess, tokenRefresh, err
	} else if err != nil {
		return tokenAccess, tokenRefresh, err
	}

	tokenAccess, err = jwt.GenerateAccessToken(s.cfg.App.SecretKey, identity, user)
	if err != nil {
		return tokenAccess, tokenRefresh, errors.New("error to generate access token: " + err.Error())
	}

	tokenRefresh, err = jwt.GenerateRefreshToken(s.cfg.App.SecretKey, identity, user)
	if err != nil {
		return tokenAccess, tokenRefresh, errors.New("error to generate refresh token: " + err.Error())
	}

	tokenRefresh.UserID = identity.UserID
	err = s.db.Debug().Model(&entity.Token{}).Create(&tokenRefresh).Error
	if err != nil {
		return
	}
	return
}
