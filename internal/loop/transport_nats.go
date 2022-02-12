package loop

import "context"

func (s service) NatsNewTransport() {
	s.DeleteRefreshTokens(context.TODO())
}

func (s service) NatsAdd(ctx context.Context) {
	l := s.logger
	l.Info("")
}
