package broker

import (
	"github.com/elga-io/corgi/internal/config"
	"github.com/elga-io/corgi/pkg/log"
	"github.com/nats-io/nats.go"
	"os"
	"time"
)

// NewBroker create a broker client object. Using NATS here.
func NewBroker(logger log.Logger, c config.Config) (broker *nats.Conn) {
	return InitBroker(logger, c)
}

// InitBroker create a broker connection.
func InitBroker(logger log.Logger, c config.Config) (broker *nats.Conn) {
	// connString := fmt.Sprintf("postgres://%s@%s:%d/%s", c.Database.User, c.Database.Host, c.Database.Port, c.Database.Base)
	nc, err := nats.Connect(nats.DefaultURL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		logger.Error("error configuring broker", "method", "InitBroker", "err", err.Error())
		os.Exit(1)
	}
	return nc
}
