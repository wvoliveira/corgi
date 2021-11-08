package server

import (
	"os"

	"github.com/go-kit/log"
)

func InitLogger() (logger log.Logger) {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return
}
