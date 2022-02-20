package redirect

import (
	"context"
)

func (s service) NATSNewTransport() {
	s.NATSFindByKeyword(context.TODO())
}

func (s service) NATSFindByKeyword(ctx context.Context) {
	l := s.logger

	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.findbykeyword", func(subj, reply string, d findByKeywordRequest) {
		l.With(ctx, "domain", d.Domain, "keyword", d.Keyword).Info("a new message received in NatsFindByKeyword")

		// Business logic.
		link, err := s.FindByKeyword(ctx, d.Domain, d.Keyword)
		if err != nil {
			l.With(ctx, "method", "s.FindByID").Error(err)
		}

		// Response.
		payload := findByKeywordResponse{Link: link, Err: err}
		err = s.broker.Publish(reply, payload)
		if err != nil {
			l.With(ctx, "method", "broker.Publish").Error(err)
		}
	})
	if err != nil {
		l.With(ctx, "method", "broker.Subscribe").Error()
	}
}
