package auth

import (
	"context"
	"time"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware middleware for all services
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostSignup(ctx context.Context, a Auth) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostAuth", "id", a.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostAuth(ctx, u)
}
