package health

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Health(ctx context.Context) ([]entity.Health, error)

	NewHTTP(r *mux.Router)
	HTTPHealth(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db      *gorm.DB
	cfg     config.Config
	version string
}

// NewService creates a new healthcheck service.
func NewService(db *gorm.DB, cfg config.Config, version string) Service {
	return service{db, cfg, version}
}

// Health make a healt check for system dependencies
// like database, social network authentication, etc.
func (s service) Health(ctx context.Context) (hs []entity.Health, err error) {
	var wg sync.WaitGroup

	// Increase async group and check a dependencie.
	// It's useful for non-blocking healthcheck.
	wg.Add(1)
	go func() {
		defer wg.Done()
		hs = append(hs, s.healthDatabase(ctx))
	}()

	// Checking Google authentication.
	wg.Add(1)
	go func() {
		defer wg.Done()
		hs = append(hs, s.healthAuthentication(ctx, "google"))
	}()

	// Checking Facebook authentication.
	wg.Add(1)
	go func() {
		defer wg.Done()
		hs = append(hs, s.healthAuthentication(ctx, "facebook"))
	}()

	wg.Wait()
	return
}

func (s service) healthDatabase(ctx context.Context) (h entity.Health) {
	h = entity.Health{
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

func (s service) healthAuthentication(ctx context.Context, provider string) (h entity.Health) {
	provider = strings.ToLower(provider)
	component := fmt.Sprintf("%s Auth", cases.Title(language.English, cases.Compact).String(provider))

	// The default config is disabled social authentication
	h = entity.Health{
		Required:    false,
		Status:      "Disabled",
		Component:   component,
		Description: fmt.Sprintf("%s config is not proper configured", component),
	}

	var endpoint string

	switch provider {
	case "google":
		endpoint = google.Endpoint.AuthURL
		if s.cfg.Auth.Google.ClientID == "" || s.cfg.Auth.Google.ClientSecret == "" {
			return
		}
	case "facebook":
		endpoint = facebook.Endpoint.AuthURL
		if s.cfg.Auth.Facebook.ClientID == "" || s.cfg.Auth.Facebook.ClientSecret == "" {
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
