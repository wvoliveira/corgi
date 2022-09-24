package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the cronjobs.
type Service interface {
	Start()
	Stop()

	RemoveTokens()
}

type service struct {
	cronn *cron.Cron
	db    *gorm.DB
	cfg   config.Config
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, cfg config.Config) Service {
	cronn := cron.New()
	return service{cronn, db, cfg}
}

func (s service) Start() {
	l := logger.Logger(context.TODO())

	err := s.cronn.AddFunc("@every 30m", func() { s.RemoveTokens() })
	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	s.cronn.Start()
	l.Info().Caller().Msg("jobs started")
}

func (s service) Stop() {
	s.cronn.Stop()
}

// RemoveTokens remove expired tokens from database.
func (s service) RemoveTokens() {
	l := logger.Logger(context.TODO())

	tokens := []entity.Token{}
	now := time.Now()

	stat := s.db.Where("? > expires_in", now).Delete(&tokens)
	count := stat.RowsAffected
	err := stat.Error

	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	l.Info().Caller().Msg(fmt.Sprintf("%d expired tokens was removed", count))
}
