package server

import (
	"os"

	"github.com/go-kit/log"
)

// NewLogger initialize a new logging object.
func NewLogger() (logger log.Logger) {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return
}
