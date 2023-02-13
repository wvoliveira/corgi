package health

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/model"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Health(*gin.Context) ([]model.Health, error)

	NewHTTP(*gin.RouterGroup)
	HTTPHealth(*gin.Context)
}

type service struct {
	db      *sql.DB
	cache   *redis.Client
	version string
}

// NewService creates a new healthcheck service.
func NewService(db *sql.DB, cache *redis.Client, version string) Service {
	return service{db, cache, version}
}

// Health make a health check for system dependencies
// like database, social network authentication, etc.
func (s service) Health(c *gin.Context) (hs []model.Health, err error) {

	hs = append(hs, s.healthDatabase(c))
	hs = append(hs, s.healthCache(c))
	hs = append(hs, s.healthAuthentication(c, "google"))
	hs = append(hs, s.healthAuthentication(c, "facebook"))

	return
}

func (s service) healthDatabase(c *gin.Context) (h model.Health) {
	h = model.Health{
		Required:    true,
		Status:      "ok",
		Component:   "database",
		Description: "integrity check is ok",
	}

	// You can mock a error with below one or create one yourself.
	// err := errors.New("database is locked. Call your admin hero")
	_, err := s.db.Exec("SELECT 1")

	if err != nil {
		h.Status = "error"
		h.Description = fmt.Sprintf("Integrity check return error: %s", err.Error())
	}

	return
}

func (s service) healthCache(c *gin.Context) (h model.Health) {
	h = model.Health{
		Required:    false,
		Status:      "ok",
		Component:   "cache",
		Description: "ping check is ok",
	}

	status := s.cache.Ping(c)

	if status.Err() != nil {
		h.Status = "error"
		h.Description = fmt.Sprintf("ping check error: %s", status.Err().Error())
	}

	return
}

func (s service) healthAuthentication(ctx context.Context, provider string) (h model.Health) {
	provider = strings.ToLower(provider)
	component := fmt.Sprintf("%s Auth", cases.Title(language.English, cases.Compact).String(provider))

	// The default config is disabled social authentication
	h = model.Health{
		Required:    false,
		Status:      "disabled",
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
		h.Status = "error"
		h.Description = err.Error()
		return
	}

	// Facebook oauth endpoint returns 500, but this not means that unavailable.
	if resp.StatusCode >= 501 && resp.StatusCode <= 599 {
		h.Status = "error"
		h.Description = fmt.Sprintf("status code from endpoint: %d", resp.StatusCode)
	} else {
		h.Status = "ok"
		h.Description = fmt.Sprintf("%s is OK", component)
	}

	return
}
