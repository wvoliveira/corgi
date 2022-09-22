package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/wvoliveira/corgi/internal/pkg/config"
	"github.com/wvoliveira/corgi/internal/pkg/entity"
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
	cfg   config.Config
}

// NewService creates a new authentication service.
func NewService(db *gorm.DB, cfg config.Config) Service {
	cronn := cron.New()
	return service{cronn, db, cfg}
}

func (s service) Start() {
	s.cronn.AddFunc("@every 30m", func() { s.RemoveTokens(context.TODO()) })
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
