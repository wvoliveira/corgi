package redirect

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/elga-io/corgi/pkg/queue"
	"log"
	"sync"
)

type ConsumerType string

const (
	// Consumers to consume messages one by one.
	// A single goroutine handles all messages.
	// Progression is slower and requires less system resource.
	// Ideal for quiet/non-critical queues.
	SyncConsumer ConsumerType = "blocking"
	// Consumers to consume messages at the same time.
	// Runs an individual goroutine per message.
	// Progression is faster and requires more system resource.
	// Ideal for busy/critical queues.
	AsyncConsumer ConsumerType = "non-blocking"
)

type ConsumerConfig struct {
	// Instructs whether to consume messages come from a worker synchronously or asynchronous.
	Type ConsumerType
	// Queue URL to receive messages from.
	QueueURL string
	// Maximum workers that will independently receive messages from a queue.
	MaxWorker int
	// Maximum messages that will be picked up by a worker in one-go.
	MaxMsg int
}

type Consumer struct {
	client queue.MessageClient
	config ConsumerConfig
}

func (s service) MQNewTransport(client queue.MessageClient, config ConsumerConfig) Consumer {
	return Consumer{
		client: client,
		config: config,
	}
}

func (c Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	for i := 1; i <= 1; i++ {
		go c.worker(ctx, wg, i)
	}

	wg.Wait()
}

func (c Consumer) worker(ctx context.Context, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	log.Printf("worker %d: started\n", id)

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker %d: stopped\n", id)
			return
		default:
		}

		fmt.Println("passei no for")

		msgs, err := c.client.Receive(ctx, "http://localhost:9324/queue/default", 10)
		if err != nil {
			// Critical error!
			log.Printf("worker %d: receive error: %s\n", id, err.Error())
			continue
		}

		if len(msgs) == 0 {
			continue
		}

		if c.config.Type == SyncConsumer {
			c.sync(ctx, msgs)
		} else {
			c.async(ctx, msgs)
		}
	}
}

func (c Consumer) sync(ctx context.Context, msgs []types.Message) {
	for _, msg := range msgs {
		c.consume(ctx, msg)
	}
}

func (c Consumer) async(ctx context.Context, msgs []types.Message) {
	wg := &sync.WaitGroup{}
	wg.Add(len(msgs))

	for _, msg := range msgs {
		go func(msg types.Message) {
			defer wg.Done()

			c.consume(ctx, msg)
		}(msg)
	}

	wg.Wait()
}

func (c Consumer) consume(ctx context.Context, msg types.Message) {
	log.Println(*msg.Body)

	if err := c.client.Delete(ctx, "http://localhost:9324/queue/default", *msg.ReceiptHandle); err != nil {
		// Critical error!
		log.Printf("delete error: %s\n", err.Error())
	}
}
