package server

import (
	"context"
	"net/http"
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

// AccessControl set common headers for web UI.
func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func (mw loggingMiddleware) SignIn(ctx context.Context, p Account) (_ Account, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "SignInAccount", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.SignIn(ctx, p)
}

func (mw loggingMiddleware) SignUp(ctx context.Context, p Account) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "SignUpAccount", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.SignUp(ctx, p)
}

func (mw loggingMiddleware) AddAccount(ctx context.Context, p Account) (_ Account, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostAccount", "id", p.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.AddAccount(ctx, p)
}

func (mw loggingMiddleware) FindAccountByID(ctx context.Context, id string) (p Account, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAccount", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.FindAccountByID(ctx, id)
}

func (mw loggingMiddleware) FindAccounts(ctx context.Context, offset, pageSize int) (p []Account, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetAccounts", "offset", offset, "page_size", pageSize, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.FindAccounts(ctx, offset, pageSize)
}

func (mw loggingMiddleware) UpdateOrCreateAccount(ctx context.Context, id string, p Account) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutAccount", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UpdateOrCreateAccount(ctx, id, p)
}

func (mw loggingMiddleware) UpdateAccount(ctx context.Context, id string, p Account) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchAccount", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UpdateAccount(ctx, id, p)
}

func (mw loggingMiddleware) DeleteAccount(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteAccount", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteAccount(ctx, id)
}

func (mw loggingMiddleware) AddURL(ctx context.Context, u URL) (_ URL, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostURL", "id", u.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.AddURL(ctx, u)
}

func (mw loggingMiddleware) FindURLByID(ctx context.Context, id string) (p URL, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.FindURLByID(ctx, id)
}

func (mw loggingMiddleware) FindURLs(ctx context.Context, offset, pageSize int) (p []URL, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetURLs", "offset", offset, "page_size", pageSize, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.FindURLs(ctx, offset, pageSize)
}

func (mw loggingMiddleware) UpdateOrCreateURL(ctx context.Context, id string, p URL) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UpdateOrCreateURL(ctx, id, p)
}

func (mw loggingMiddleware) UpdateURL(ctx context.Context, id string, p URL) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.UpdateURL(ctx, id, p)
}

func (mw loggingMiddleware) DeleteURL(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteURL", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteURL(ctx, id)
}
