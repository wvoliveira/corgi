package logger

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

func Logger(ctx context.Context) (l zerolog.Logger) {
	logContext := log.Logger.With()

	if ctxRequest, ok := ctx.Value(model.CorrelationID{}).(model.CorrelationID); ok {
		logContext = logContext.Str("req_id", ctxRequest.ID)
	}

	if ctxIdentity, ok := ctx.Value(model.IdentityInfo{}).(model.IdentityInfo); ok {
		logContext = logContext.Str("user_id", ctxIdentity.UserID)
	}

	return logContext.Logger()
}
