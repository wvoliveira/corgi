package url

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

func (s *instrumentingService) PostURL(ctx context.Context, u URL) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "new").Add(1)
		s.requestLatency.With("method", "new").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.PostURL(ctx, u)
}

func (s *instrumentingService) GetURL(ctx context.Context, id string) (URL, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetURL(ctx, id)
}

func (s *instrumentingService) GetURLs(ctx context.Context, offset, pageSize int) ([]URL, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list_urls").Add(1)
		s.requestLatency.With("method", "list_urls").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetURLs(ctx, offset, pageSize)
}
