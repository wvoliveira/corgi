package auth

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

func (s *instrumentingService) PostSignup(ctx context.Context, a Auth) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "signup").Add(1)
		s.requestLatency.With("method", "signup").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostSignuup(ctx, a)
}
