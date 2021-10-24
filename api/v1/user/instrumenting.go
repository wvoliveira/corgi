package user

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) PostUser(ctx context.Context, u User) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "new").Add(1)
		s.requestLatency.With("method", "new").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostUser(ctx, u)
}

func (s *instrumentingService) GetUser(ctx context.Context, id string) (User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetUser(ctx, id)
}

func (s *instrumentingService) GetUsers(ctx context.Context, offset, pageSize int) ([]User, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list_users").Add(1)
		s.requestLatency.With("method", "list_users").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetUsers(ctx, offset, pageSize)
}
