package pwd

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

func (mw loggingMiddleware) SignInPwd(ctx context.Context, p Pwd) (_ Pwd, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "SignInPwd", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.SignInPwd(ctx, p)
}

func (mw loggingMiddleware) SignUpPwd(ctx context.Context, p Pwd) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "SignUpPwd", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.SignUpPwd(ctx, p)
}
