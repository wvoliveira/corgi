package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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

func Default() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	level := GetLogLevel()
	zerolog.SetGlobalLevel(level)
	// gin.SetMode(gin.ReleaseMode)
}

func GetLogLevel() zerolog.Level {
	levels := map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
	}

	configLogLevel := viper.GetString("app.log_level")
	level, exists := levels[configLogLevel]
	if exists {
		return level
	}

	return zerolog.InfoLevel
}
