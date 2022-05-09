package logger

import (
	"context"

	"github.com/elga-io/corgi/internal/app/entity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(ctx context.Context) (l zerolog.Logger) {
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(entity.CorrelationID{}).(entity.CorrelationID); ok {
			l = log.With().Str("req_id", ctxRqId.ID).Logger()
		}
	}
	return l
}
