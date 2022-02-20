package queue

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type MessageClient interface {
	// Send a message to queue and returns its message ID.
	Send(ctx context.Context, req *SendRequest) (string, error)
	// Long polls given amount of messages from a queue.
	Receive(ctx context.Context, queueURL string, maxMsg int32) ([]types.Message, error)
	// Deletes a message from a queue.
	Delete(ctx context.Context, queueURL, rcvHandle string) error
}

type SendRequest struct {
	QueueURL   string
	Body       string
	Attributes []Attribute
}

type Attribute struct {
	Key   string
	Value string
	Type  string
}
