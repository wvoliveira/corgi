package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
	"github.com/wvoliveira/corgi/internal/pkg/logger"
	"gorm.io/gorm"
)

// Service encapsulates the cronjobs.
type Service interface {
	Start()
	Stop()

	RemoveTokens(ctx context.Context)
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
	l := logger.Logger(context.TODO())

	err := s.cronn.AddFunc("@every 30m", func() { s.RemoveTokens(context.TODO()) })
	if err != nil {
		l.Error().Caller().Msg(err.Error())
	}

	s.cronn.Start()
}

func (s service) Stop() {
	s.cronn.Stop()
}

// RemoveTokens remove expired tokens from database.
func (s service) RemoveTokens(_ context.Context) {
	tokens := []entity.Token{}
	now := time.Now()

	err := s.db.Where("? > expires_in", now).Delete(&tokens).Error
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
