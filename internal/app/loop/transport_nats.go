package loop

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (s service) NatsNewTransport() {
	s.DeleteRefreshTokens(context.TODO())
}

func (s service) NatsAdd(ctx context.Context) {
	l := log.Ctx(ctx)
	l.Info().Caller().Discard().Send()
}
