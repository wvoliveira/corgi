package google

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	e "github.com/wvoliveira/corgi/internal/pkg/errors"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// Service encapsulates the authentication logic.
type Service interface {
	Login(*gin.Context, string, string) (model.User, string, error)
	Callback(*gin.Context, string, callbackRequest) (model.User, error)

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
func (s service) Login(c *gin.Context, accessToken, callbackURL string) (user model.User, redirectURL string, err error) {
	log := logger.Logger(c)

	conf := createOAuth2Config(callbackURL)

	if accessToken == "" {
		// Redirect user to Google's consent page to ask for permission
		// for the scopes specified above.
		redirectURL = conf.AuthCodeURL("state")
		return
	}

	// If access token was sent with URL, just get user info.
	userGoogle, err := getUserFromGoogle(c, accessToken)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, user, err = getOrCreateUser(s.db, userGoogle)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

func (s service) Callback(c *gin.Context, callbackURL string, r callbackRequest) (user model.User, err error) {
	log := logger.Logger(c)

	conf := createOAuth2Config(callbackURL)

	oauthToken, err := conf.Exchange(c, r.Code)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	userGoogle, err := getUserFromGoogle(c, oauthToken.AccessToken)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, user, err = getOrCreateUser(s.db, userGoogle)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

func createOAuth2Config(callbackURL string) (conf *oauth2.Config) {
	conf = &oauth2.Config{
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
	return
}

func getUserFromGoogle(c *gin.Context, accessToken string) (userGoogle model.UserGoogle, err error) {

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(accessToken))

	if err != nil {
		return
	}

	if resp.StatusCode == 401 {
		return userGoogle, e.ErrUnauthorized
	}

	if resp.StatusCode != 200 {
		return userGoogle, errors.New("error to get info from Google")
	}

	if resp.StatusCode == 401 {
		return userGoogle, e.ErrUnauthorized
	}

	if resp.StatusCode != 200 {
		return userGoogle, errors.New("error to get info from Google")
	}

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &userGoogle)
	return
}

func getOrCreateUser(db *gorm.DB, userGoogle model.UserGoogle) (identity model.Identity, user model.User, err error) {

	err = db.
		Model(model.Identity{}).
		Where("provider = ? AND UID = ?", "google", userGoogle.ID).
		First(&identity).Error

	if err == gorm.ErrRecordNotFound {

		identity.ID = uuid.New().String()
		identity.CreatedAt = time.Now()
		identity.LastLogin = identity.CreatedAt
		identity.Provider = "google"
		identity.UID = userGoogle.ID
		identity.Verified = &userGoogle.VerifiedEmail

		active := true
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.Name = userGoogle.Name
		user.Role = "user"
		user.Active = &active
		user.Identities = append(user.Identities, identity)

		err = db.
			Model(&model.User{}).
			Create(&user).Error

		return

	}

	if err != nil {
		return
	}

	if identity.UserID == "" {
		err = db.
			Model(&model.Identity{}).
			Where("provider = ? AND uid = ?", "google", userGoogle.ID).
			First(&identity).Error

		if err != nil {
			return
		}
	}

	err = db.
		Model(&model.User{}).
		Where("id = ?", identity.UserID).
		First(&user).Error

	return
}
