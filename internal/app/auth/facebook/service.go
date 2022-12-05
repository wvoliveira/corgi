package facebook

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
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
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
		// Redirect user to Facebook's consent page to ask for permission
		// for the scopes specified above.
		redirectURL = conf.AuthCodeURL("state")
		return
	}

	// If access token was sent with URL, just get user info.
	userFacebook, err := getUserFromFacebook(c, accessToken)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, user, err = getOrCreateUser(c, s.db, userFacebook)

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

	userFacebook, err := getUserFromFacebook(c, oauthToken.AccessToken)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	_, user, err = getOrCreateUser(c, s.db, userFacebook)

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	return
}

func createOAuth2Config(callbackURL string) (conf *oauth2.Config) {
	conf = &oauth2.Config{
		ClientID:     viper.GetString("auth.facebook.client_id"),
		ClientSecret: viper.GetString("auth.facebook.client_secret"),
		RedirectURL:  callbackURL,
		Scopes: []string{
			"public_profile",
			"email",
		},
		Endpoint: facebook.Endpoint,
	}
	return
}

func getUserFromFacebook(c *gin.Context, accessToken string) (userFacebook model.UserFacebook, err error) {

	resp, err := http.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + url.QueryEscape(accessToken))

	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return userFacebook, errors.New("error to get info from Facebook")
	}

	response, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(response, &userFacebook)
	return
}

func getOrCreateUser(c *gin.Context, db *gorm.DB, userFacebook model.UserFacebook) (identity model.Identity, user model.User, err error) {
	log := logger.Logger(c)

	err = db.Debug().Model(model.Identity{}).Where("provider = ? AND UID = ?", "facebook", userFacebook.ID).First(&identity).Error

	if err == gorm.ErrRecordNotFound {

		identity.ID = uuid.New().String()
		identity.CreatedAt = time.Now()
		identity.LastLogin = identity.CreatedAt
		identity.Provider = "facebook"
		identity.UID = userFacebook.ID

		active := true
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.Name = userFacebook.Name
		user.Role = "user"
		user.Active = &active
		user.Identities = append(user.Identities, identity)

		err = db.Debug().Model(&model.User{}).Create(&user).Error
		return

	}

	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return
	}

	if identity.UserID == "" {
		err = db.Debug().Model(&model.Identity{}).Where("provider = ? AND uid = ?", "facebook", userFacebook.ID).First(&identity).Error

		if err != nil {
			return
		}
	}

	err = db.Debug().Model(&model.User{}).Where("id = ?", identity.UserID).First(&user).Error

	return
}
