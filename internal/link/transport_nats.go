package link

import (
	"context"
	"github.com/elga-io/corgi/internal/entity"
)

func (s service) NatsNewTransport() {
	s.NatsAdd(context.TODO())
	s.NatsFindByID(context.TODO())
	s.NatsFindAll(context.TODO())
	s.NatsUpdate(context.TODO())
	s.NatsDelete(context.TODO())
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
		payload := addResponse{ID: link.ID, Domain: link.Domain, Keyword: link.Keyword, URL: link.URL, Title: link.Title, Err: err}
		err = s.broker.Publish(reply, payload)
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
		payload := findByIDResponse{Link: link, Err: err}
		err = s.broker.Publish(reply, payload)
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
		payload := findAllResponse{Links: links, Limit: d.Limit, Page: d.Page, Sort: d.Sort, Total: total, Pages: pages, Err: err}
		err = s.broker.Publish(reply, payload)
		if err != nil {
			l.With(ctx, "method", "broker.Publish").Error(err)
		}
	})
	if err != nil {
		l.With(ctx, "method", "broker.Subscribe").Error()
	}
}

func (s service) NatsUpdate(ctx context.Context) {
	l := s.logger

	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.update", func(subj, reply string, d updateRequest) {
		l.With(ctx, "id", d.ID, "user_id", d.UserID).
			Info("a new message received in NatsUpdate")

		// Business logic.
		data := entity.Link{ID: d.ID, Domain: d.Domain, Keyword: d.Keyword, URL: d.URL, Title: d.Title, Active: d.Active, UserID: d.UserID}

		link, err := s.Update(ctx, data)
		if err != nil {
			l.With(ctx, "method", "NatsUpdate").Error(err)
		}

		// Response.
		payload := updateResponse{Link: link, Err: err}
		err = s.broker.Publish(reply, payload)
		if err != nil {
			l.With(ctx, "method", "broker.Publish").Error(err)
		}
	})
	if err != nil {
		l.With(ctx, "method", "broker.Subscribe").Error()
	}
}

func (s service) NatsDelete(ctx context.Context) {
	l := s.logger

	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.delete", func(subj, reply string, d deleteRequest) {
		l.With(ctx, "id", d.ID, "user_id", d.UserID).
			Info("a new message received in NatsDelete")

		// Business logic.
		err := s.Delete(ctx, d.ID, d.UserID)
		if err != nil {
			l.With(ctx, "method", "NatsDelete").Error(err)
		}

		// Response.
		payload := deleteResponse{Err: err}
		err = s.broker.Publish(reply, payload)
		if err != nil {
			l.With(ctx, "method", "broker.Publish").Error(err)
		}
	})
	if err != nil {
		l.With(ctx, "method", "broker.Subscribe").Error()
	}
}
