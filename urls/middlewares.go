package urls

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

func (mw loggingMiddleware) PostURL(ctx context.Context, p URL) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostURL", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostURL(ctx, p)
}

func (mw loggingMiddleware) GetURL(ctx context.Context, id string) (p URL, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetURL(ctx, id)
}

func (mw loggingMiddleware) GetURLs(ctx context.Context, offset, pageSize int) (p []URL, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetURLs", "offset", offset, "page_size", pageSize, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetURLs(ctx, offset, pageSize)
}

func (mw loggingMiddleware) PutURL(ctx context.Context, id string, p URL) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutURL(ctx, id, p)
}

func (mw loggingMiddleware) PatchURL(ctx context.Context, id string, p URL) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchURL(ctx, id, p)
}

func (mw loggingMiddleware) DeleteURL(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteURL(ctx, id)
}
