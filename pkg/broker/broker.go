package broker

import (
	"context"
	"fmt"
	"github.com/elga-io/corgi/internal/config"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

// NewBroker create a broker client object. Using NATS here.
func NewBroker(logger log.Logger, c config.Config) (broker *nats.EncodedConn) {
	return InitBroker(logger, c)
}

// InitBroker create a broker connection.
func InitBroker(logger log.Logger, c config.Config) (broker *nats.EncodedConn) {
	connURL := fmt.Sprintf("nats://%v:%v", c.Broker.Host, c.Broker.Port)
	nc, err := nats.Connect(connURL,
		nats.UserInfo("foo", "bar"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		logger.With(context.TODO(), "err", err.Error()).Error("error configuring broker")
		os.Exit(1)
	}

	conn, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		logger.With(context.TODO(), "err", err.Error()).Error("error configuring new encoded connection")
		os.Exit(1)
	}

	return conn
}
