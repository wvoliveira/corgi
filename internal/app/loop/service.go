package loop

/*
This package serves functions that need to run frequently. Like:
- Delete expired tokens
- Renew subscription
*/

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Service encapsulates the link service logic, http handlers and another transport layer.
type Service interface {
	DeleteRefreshTokens(ctx context.Context)
}

type service struct {
	db     *gorm.DB
	broker *nats.EncodedConn
}

// NewService creates a new loop service.
func NewService(db *gorm.DB, broker *nats.EncodedConn) Service {
	return service{db, broker}
}

func (s service) DeleteRefreshTokens(ctx context.Context) {
	l := log.Ctx(ctx)
	l.Info().Caller().Msg("Init delete refresh token routine")
}
