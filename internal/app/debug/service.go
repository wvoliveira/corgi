package debug

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	Info(ctx context.Context) (info, error)

	NewHTTP(r *mux.Router)
	HTTPInfo(w http.ResponseWriter, r *http.Request)
}

type service struct {
	db      *gorm.DB
	cfg     config.Config
	version string
}

type info struct {
	Version     string     `json:"version"`
	LogLevel    string     `json:"log_level"`
	RedirectURL string     `json:"redirect_url"`
	Server      infoServer `json:"server"`
}

type infoServer struct {
	HTTPPort string `json:"http_port"`
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, cfg config.Config, version string) Service {
	return service{db, cfg, version}
}

// Health create a new shortener link.
func (s service) Info(_ context.Context) (i info, err error) {
	// l := logger.Logger(ctx)
	i.Version = s.version
	i.LogLevel = s.cfg.LogLevel
	i.RedirectURL = s.cfg.RedirectURL
	i.Server.HTTPPort = s.cfg.Server.HTTPPort
	return
}
