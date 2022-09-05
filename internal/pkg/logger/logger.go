package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/app/entity"
)

func Logger(ctx context.Context) (l zerolog.Logger) {
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(entity.CorrelationID{}).(entity.CorrelationID); ok {
			l = log.With().Str("req_id", ctxRqId.ID).Logger()
		}
	}
	return l
}
