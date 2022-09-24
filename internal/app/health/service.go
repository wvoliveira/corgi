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
	HealthDatabase(ctx context.Context) (h entity.Health)
	HealthAuth(ctx context.Context, providers []string) ([]entity.Health, error)

	NewHTTP(r *mux.Router)
	HTTPHealth(w http.ResponseWriter, r *http.Request)
	HTTPHealthAuth(w http.ResponseWriter, r *http.Request)
	HTTPHealthAuthProvider(w http.ResponseWriter, r *http.Request)
	HTTPHealthLive(w http.ResponseWriter, r *http.Request)
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

// Health make a health check for system dependencies
// like database, social network authentication, etc.
func (s service) Health(ctx context.Context) (hs []entity.Health, err error) {
	// Increase async group and check a dependence.
	// It's useful for non-blocking healthcheck.
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		hs = append(hs, s.HealthDatabase(ctx))
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

// HealthAuth make a health check for auth dependencies.
// like Google auth, Facebook auth or another social network.
func (s service) HealthAuth(ctx context.Context, providers []string) (hs []entity.Health, err error) {
	// Increase async group and check a dependence.
	// It's useful for non-blocking healthcheck.
	var wg sync.WaitGroup

	if len(providers) == 0 {
		providers = []string{"google", "facebook"}
	}

	for index := range providers {
		wg.Add(1)

		// God knows what this means
		// and today I too, but who knows in some days.
		go func(provider string) {
			defer wg.Done()
			hs = append(hs, s.healthAuthentication(ctx, provider))
		}(providers[index])

	}

	wg.Wait()
	return
}

func (s service) HealthDatabase(ctx context.Context) (h entity.Health) {
	h = entity.Health{
		Required:    true,
		Status:      "ok",
		Component:   "database",
		Description: "integrity check is ok",
	}

	// You can mock a error with below one or create one yourself.
	// err := errors.New("database is locked. Call your admin hero")
	err := s.db.Exec("PRAGMA integrity_check").Error
	if err != nil {
		h.Status = "error"
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
		Status:      "disabled",
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
