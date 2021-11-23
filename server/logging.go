package server

import (
	"go.uber.org/zap"
)

// NewLogger initialize a new logging object.
func NewLogger() (logger *zap.SugaredLogger) {
	zlog, _ := zap.NewProduction()
	_ = zlog.Sync() // flushes buffer, if any
	return zlog.Sugar()
}
