package server

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

// NewMetrics start all metrics services.
func NewMetrics(s Service) Service {
	s = newMetricsAuth(s)
	s = newMetricsAccount(s)
	s = newMetricsURL(s)
	return s
}

/*
	Struct for instrument Auth.
*/
type instrumentingAuth struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

/*
	Create an instance of an instrumenting Auth service.
*/
func newMetricsAuth(s Service) Service {
	counter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "auth_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{})

	latency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "api",
		Subsystem: "auth_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, []string{})

	return &instrumentingAuth{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingAuth) SignIn(ctx context.Context, u Account) (Account, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "sign_in").Add(1)
		s.requestLatency.With("method", "sign_in").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.SignIn(ctx, u)
}

func (s *instrumentingAuth) SignUp(ctx context.Context, u Account) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "sign_up").Add(1)
		s.requestLatency.With("method", "sign_up").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.SignUp(ctx, u)
}

/*
	Struct for instrument Account.
*/
type instrumentingAccount struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

/*
	Create an instance of an instrumenting Account service.
*/
func newMetricsAccount(s Service) Service {
	counter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "account_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{})

	latency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "api",
		Subsystem: "account_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, []string{})

	return &instrumentingAccount{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingAccount) AddAccount(ctx context.Context, u Account) (Account, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "add_account").Add(1)
		s.requestLatency.With("method", "add_account").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.AddAccount(ctx, u)
}

func (s *instrumentingAccount) FindAccountByID(ctx context.Context, id string) (Account, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find_account_by_id").Add(1)
		s.requestLatency.With("method", "find_account_by_id").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.FindAccountByID(ctx, id)
}

func (s *instrumentingAccount) FindAccounts(ctx context.Context, offset, pageSize int) ([]Account, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find_accounts").Add(1)
		s.requestLatency.With("method", "find_accounts").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.FindAccounts(ctx, offset, pageSize)
}

func (s *instrumentingAccount) UpdateOrCreateAccount(ctx context.Context, id string, reqAccount Account) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find_or_create_account").Add(1)
		s.requestLatency.With("method", "find_or_create_account").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.UpdateOrCreateAccount(ctx, id, reqAccount)
}

func (s *instrumentingAccount) UpdateAccount(ctx context.Context, id string, reqAccount Account) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "update_account").Add(1)
		s.requestLatency.With("method", "update_account").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.UpdateAccount(ctx, id, reqAccount)
}

func (s *instrumentingAccount) DeleteAccount(ctx context.Context, id string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "update_account").Add(1)
		s.requestLatency.With("method", "update_account").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.DeleteAccount(ctx, id)
}

/*
	Struct for instrument URL.
*/
type instrumentingURL struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

/*
	Create an instance of an instrumenting URL service.
*/
func newMetricsURL(s Service) Service {
	counter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "url_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{})

	latency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "api",
		Subsystem: "url_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, []string{})

	return &instrumentingURL{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingURL) AddURL(ctx context.Context, u URL) (URL, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "add_url").Add(1)
		s.requestLatency.With("method", "add_url").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.AddURL(ctx, u)
}

func (s *instrumentingURL) FindURLbyID(ctx context.Context, id string) (URL, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.FindURLByID(ctx, id)
}

func (s *instrumentingURL) FindURLs(ctx context.Context, offset, pageSize int) ([]URL, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list_urls").Add(1)
		s.requestLatency.With("method", "list_urls").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.FindURLs(ctx, offset, pageSize)
}

func (s *instrumentingURL) UpdateOrCreateURL(ctx context.Context, id string, reqURL URL) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "find_or_create_url").Add(1)
		s.requestLatency.With("method", "find_or_create_url").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.UpdateOrCreateURL(ctx, id, reqURL)
}

func (s *instrumentingURL) UpdateURL(ctx context.Context, id string, reqURL URL) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "update_url").Add(1)
		s.requestLatency.With("method", "update_url").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.UpdateURL(ctx, id, reqURL)
}

func (s *instrumentingURL) DeleteURL(ctx context.Context, id string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "update_url").Add(1)
		s.requestLatency.With("method", "update_url").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.DeleteURL(ctx, id)
}
