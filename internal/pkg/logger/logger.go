package logger

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/wvoliveira/corgi/internal/pkg/model"
)

func Logger(c context.Context) (l zerolog.Logger) {
	logContext := log.Logger.With()

	if ctxRequest, ok := c.Value(model.CorrelationID{}).(model.CorrelationID); ok {
		logContext = logContext.Str("req_id", ctxRequest.ID)
	}

	if user_id, ok := c.Value("user_id").(string); ok {
		logContext = logContext.Str("user_id", user_id)
	}

	return logContext.Logger()
}

func Default() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	level := GetLogLevel()
	zerolog.SetGlobalLevel(level)

	if level == zerolog.InfoLevel {
		gin.SetMode(gin.ReleaseMode)
	}
}

func GetLogLevel() zerolog.Level {
	levels := map[string]zerolog.Level{
		"trace": zerolog.TraceLevel,
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
	}

	configLogLevel := viper.GetString("LOG_LEVEL")

	level, exists := levels[configLogLevel]
	if exists {
		return level
	}

	return zerolog.InfoLevel
}
