package profile

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

func (s *instrumentingService) PostProfile(ctx context.Context, u Profile) (Profile, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "new").Add(1)
		s.requestLatency.With("method", "new").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostProfile(ctx, u)
}

func (s *instrumentingService) GetProfile(ctx context.Context, id string) (Profile, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetProfile(ctx, id)
}

func (s *instrumentingService) GetProfiles(ctx context.Context, offset, pageSize int) ([]Profile, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list_profiles").Add(1)
		s.requestLatency.With("method", "list_profiles").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetProfiles(ctx, offset, pageSize)
}
