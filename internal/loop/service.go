package loop

/*
This package serves functions that need to run frequently. Like:
- Delete expired tokens
- Renew subscription
*/

import (
	"context"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	DeleteRefreshTokens(ctx context.Context)
}

type service struct {
	logger log.Logger
	db     *gorm.DB
	broker *nats.EncodedConn
}

// NewService creates a new loop service.
func NewService(logger log.Logger, db *gorm.DB, broker *nats.EncodedConn) Service {
	return service{logger, db, broker}
}

func (s service) DeleteRefreshTokens(ctx context.Context) {
	l := s.logger.With(ctx)
	l.Info("Init delete refresh token routine")
}
