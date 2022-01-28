package link

import (
	"context"
	"github.com/elga-io/corgi/internal/entity"
)

func (s service) NatsNewTransport() {
	s.NatsAdd(context.TODO())
	s.NatsFindByID(context.TODO())
	s.NatsFindAll(context.TODO())
}

func (s service) NatsAdd(ctx context.Context) {
	l := s.logger

	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.add", func(subj, reply string, d addRequest) {
		l.With(ctx, "domain", d.Domain, "keyword", d.Keyword, "url", d.URL, "title", d.Title, "user_id", d.UserID).Info()

		// Business logic.
		link, err := s.Add(ctx, entity.Link{Domain: d.Domain, Keyword: d.Keyword, URL: d.URL, Title: d.Title, UserID: d.UserID})
		if err != nil {
			l.With(ctx, "method", "s.Add").Error(err)
		}

		// Response.
		err = s.broker.Publish(reply, link)
		if err != nil {
			l.With(ctx, "method", "broker.Publish").Error(err)
		}
	})
	if err != nil {
		l.With(ctx, "method", "broker.Subscribe").Error(err)
	}
}

func (s service) NatsFindByID(ctx context.Context) {
	l := s.logger

	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.findbyid", func(subj, reply string, d findByIDRequest) {
		l.With(ctx, "id", d.ID, "user_id", d.UserID).Info("a new message received in NatsFindByID")

		// Business logic.
		link, err := s.FindByID(ctx, d.ID, d.UserID)
		if err != nil {
			l.With(ctx, "method", "s.FindByID").Error(err)
		}

		// Response.
		err = s.broker.Publish(reply, link)
		if err != nil {
			l.With(ctx, "method", "broker.Publish").Error(err)
		}
	})
	if err != nil {
		l.With(ctx, "method", "broker.Subscribe").Error()
	}
}

func (s service) NatsFindAll(ctx context.Context) {
	l := s.logger

	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.findall", func(subj, reply string, d findAllRequest) {
		l.With(ctx, "page", d.Page, "sort", d.Sort, "offset", d.Offset, "limit", d.Limit, "user_id", d.UserID).
			Info("a new message received in NatsFindAll")

		// Business logic.
		total, pages, links, err := s.FindAll(ctx, d.Offset, d.Limit, d.Sort, d.UserID)
		if err != nil {
			l.With(ctx, "method", "s.FindByID").Error(err)
		}

		// Response.
		err = s.broker.Publish(reply, findAllResponse{Links: links, Pages: pages, Total: total})
		if err != nil {
			l.With(ctx, "method", "broker.Publish").Error(err)
		}
	})
	if err != nil {
		l.With(ctx, "method", "broker.Subscribe").Error()
	}
}
