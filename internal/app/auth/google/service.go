package google

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(*gin.Context, string) (string, error)
	Callback(*gin.Context, string, callbackRequest) (entity.User, error)

	NewHTTP(*gin.RouterGroup)
	HTTPLogin(*gin.Context)
	HTTPCallback(*gin.Context)
}

type service struct {
	db *gorm.DB
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB) Service {
	return service{db}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(_ *gin.Context, callbackURL string) (redirectURL string, err error) {
	conf := &oauth2.Config{
		ClientID:     viper.GetString("auth.google.client_id"),
		ClientSecret: viper.GetString("auth.google.client_secret"),
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

func (s service) Callback(c *gin.Context, callbackURL string, r callbackRequest) (user entity.User, err error) {
	l := logger.Logger(c)

	conf := &oauth2.Config{
		ClientID:     viper.GetString("auth.google.client_id"),
		ClientSecret: viper.GetString("auth.google.client_secret"),
		RedirectURL:  callbackURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}

	oauthToken, err := conf.Exchange(c, r.Code)

	if err != nil {
		l.Error().Caller().Msg(err.Error())
		return
	}

	client := conf.Client(c, oauthToken)

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

		identity.ID = uuid.New().String()
		identity.CreatedAt = time.Now()
		identity.LastLogin = identity.CreatedAt
		identity.Provider = "google"
		identity.UID = googleUser.ID
		identity.Verified = &googleUser.VerifiedEmail
		identity.VerifiedAt = identity.CreatedAt

		active := true
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.Name = googleUser.Name
		user.Role = "user"
		user.Active = &active
		user.Identities = append(user.Identities, identity)

		err = s.db.Debug().Model(&entity.User{}).Create(&user).Error
		if err != nil {
			l.Error().Caller().Msg(err.Error())
			return
		}

	} else if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	if identity.UserID == "" {
		err = s.db.Debug().Model(&entity.Identity{}).Where("provider = ? AND uid = ?", "google", googleUser.ID).First(&identity).Error
		if err != nil {
			l.Error().Caller().Msg(err.Error())
			return
		}
	}

	err = s.db.Debug().Model(&entity.User{}).Where("id = ?", identity.UserID).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, err
	}

	return
}
