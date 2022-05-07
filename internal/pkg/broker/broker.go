package broker

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/elga-io/corgi/internal/app/config"
	"github.com/nats-io/nats.go"
	"log"
)

// NewBroker create a broker client object. Using NATS here.
func NewBroker(logger log.Logger, c config.Config) (broker *nats.EncodedConn) {
	return InitBroker(logger, c)
}

// InitBroker create a broker connection.
func InitBroker(logger log.Logger, c config.Config) (broker *nats.EncodedConn) {
	connURL := fmt.Sprintf("nats://%v:%v", c.Broker.Host, c.Broker.Port)
	nc, err := nats.Connect(connURL,
		nats.UserInfo("foo", "bar"), // TODO: change stream auth
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Warnf("Got disconnected! Reason: %s", err.Error())
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Infof("Got reconnected to %v", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Infof("Connection closed. Reason: %q", nc.LastError())
		}),
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
