package job

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"github.com/wvoliveira/corgi/internal/pkg/model"
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
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB) Service {
	cronn := cron.New()
	return service{cronn, db}
}

func (s service) Start() {
	log := logger.Logger(context.TODO())

	err := s.cronn.AddFunc("@every 30m", func() { s.RemoveTokens(context.TODO()) })

	if err != nil {
		log.Error().Caller().Msg(err.Error())
	}

	s.cronn.Start()
	l.Info().Caller().Msg("jobs started")
}

func (s service) Stop() {
	s.cronn.Stop()
}

// RemoveTokens remove expired tokens from database.
func (s service) RemoveTokens(_ context.Context) {
	tokens := []model.Token{}
	now := time.Now()

	stat := s.db.Where("? > expires_in", now).Delete(&tokens)
	count := stat.RowsAffected
	err := stat.Error

	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	l.Info().Caller().Msg(fmt.Sprintf("%d expired tokens was removed", count))
}
