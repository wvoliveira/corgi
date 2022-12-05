package health

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Health(*gin.Context) ([]model.Health, error)

	NewHTTP(*gin.RouterGroup)
	HTTPHealth(*gin.Context)
}

type service struct {
	db      *gorm.DB
	version string
}

// NewService creates a new healthcheck service.
func NewService(db *gorm.DB, version string) Service {
	return service{db, version}
}

// Health make a healt check for system dependencies
// like database, social network authentication, etc.
func (s service) Health(c *gin.Context) (hs []model.Health, err error) {

	var wg sync.WaitGroup

	// Increase async group and check a dependencie.
	// It's useful for non-blocking healthcheck.
	wg.Add(1)
	go func() {
		defer wg.Done()
		hs = append(hs, s.healthDatabase(c))
	}()

	// Checking Google authentication.
	wg.Add(1)
	go func() {
		defer wg.Done()
		hs = append(hs, s.healthAuthentication(c, "google"))
	}()

	// Checking Facebook authentication.
	wg.Add(1)
	go func() {
		defer wg.Done()
		hs = append(hs, s.healthAuthentication(c, "facebook"))
	}()

	wg.Wait()

	return
}

func (s service) healthDatabase(c *gin.Context) (h model.Health) {

	h = model.Health{
		Required:    true,
		Status:      "OK",
		Component:   "Database",
		Description: "Integrity check is OK",
	}

	// You can mock a error with below one or create one yourself.
	// err := errors.New("database is locked. Call your admin hero")
	err := s.db.Exec("PRAGMA integrity_check").Error

	if err != nil {
		h.Status = "Fail"
		h.Description = fmt.Sprintf("Integrity check return error: %s", err.Error())
	}

	return
}

func (s service) healthAuthentication(ctx context.Context, provider string) (h model.Health) {

	provider = strings.ToLower(provider)
	component := fmt.Sprintf("%s Auth", cases.Title(language.English, cases.Compact).String(provider))

	// The default config is disabled social authentication
	h = model.Health{
		Required:    false,
		Status:      "Disabled",
		Component:   component,
		Description: fmt.Sprintf("%s config is not proper configured", component),
	}

	var endpoint string

	switch provider {
	case "google":
		endpoint = google.Endpoint.AuthURL
		if viper.GetString("auth.google.client_id") == "" || viper.GetString("auth.google.client_secret") == "" {
			return
		}
	case "facebook":
		endpoint = facebook.Endpoint.AuthURL
		if viper.GetString("auth.facebook.client_id") == "" || viper.GetString("auth.facebook.client_secret") == "" {
			return
		}
	}

	resp, err := http.Head(endpoint)
	if err != nil {
		h.Status = "Failed"
		h.Description = err.Error()
		return
	}

	// Facebook oauth endpoint returns 500, but this not means that unavailable.
	if resp.StatusCode >= 501 && resp.StatusCode <= 599 {
		h.Status = "Failed"
		h.Description = fmt.Sprintf("Status code from endpoint: %d", resp.StatusCode)
	} else {
		h.Status = "OK"
		h.Description = fmt.Sprintf("%s is OK", component)
	}

	return
}
