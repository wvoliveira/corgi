package link

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
)

func (s service) NatsNewTransport() {
	s.NatsAdd(context.TODO())
	s.NatsFindByID(context.TODO())
}

func (s service) NatsAdd(ctx context.Context) {
	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.add", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))

		err := m.Respond([]byte("nats add"))
		if err != nil {
			s.logger.With(ctx, "error to respond to broker", err.Error())
		}
	})
	if err != nil {
		s.logger.With(ctx, "error to subscribe", err.Error())
	}
}

func (s service) NatsFindByID(ctx context.Context) {
	ll := s.logger.With(ctx, "method", "NatsFindByID")

	// Subscribe to channel and receive a message.
	_, err := s.broker.Subscribe("link.findbyid", func(subj, reply string, d findByIDRequest) {
		ll.Info("received a message", "id", d.ID, "user_id", d.UserID)

		// Business logic.
		link, err := s.FindByID(ctx, d.ID, d.UserID)
		if err != nil {
			s.logger.With(ctx, "error to find by id", err.Error())
		}

		// Response.
		err = s.broker.Publish(reply, link)
		if err != nil {
			s.logger.With(ctx, "error to respond to broker", err.Error())
		}
	})
	if err != nil {
		s.logger.With(ctx, "error to subscribe", err.Error())
	}
}
