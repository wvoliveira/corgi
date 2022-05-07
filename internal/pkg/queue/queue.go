package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/elga-io/corgi/internal/app/config"
	"log"
)

type Queuer struct {
	timeout time.Duration
	client  *sqs.Client
}

// NewQueuer create a queue client object. Using SQS or ElasticMQ here.
func NewQueuer(l log.Logger, c config.Config) Queuer {
	return Queuer{
		timeout: time.Second * 10,
		client:  initQueuer(l, c),
	}
}

func initQueuer(logger log.Logger, c config.Config) (queuer *sqs.Client) {
	// Config sources can be passed to LoadDefaultConfig, these sources can implement one or more
	// provider interfaces. These sources take priority over the standard environment and shared configuration values.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "http://localhost:9324",
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg := aws.Config{
		Region:                      "us-east-1",
		EndpointResolverWithOptions: customResolver,
	}
	return sqs.NewFromConfig(cfg)
}

func (s Queuer) Send(ctx context.Context, req *SendRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	attrs := make(map[string]types.MessageAttributeValue, len(req.Attributes))
	for _, attr := range req.Attributes {
		attrs[attr.Key] = types.MessageAttributeValue{
			StringValue: aws.String(attr.Value),
			DataType:    aws.String(attr.Type),
		}
	}

	res, err := s.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageAttributes: attrs,
		MessageBody:       aws.String(req.Body),
		QueueUrl:          aws.String(req.QueueURL),
	})
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	return *res.MessageId, nil
}

func (s Queuer) Receive(ctx context.Context, queueURL string, maxMsg int32) ([]types.Message, error) {
	if maxMsg < 1 || maxMsg > 10 {
		return nil, fmt.Errorf("receive argument: msgMax valid values: 1 to 10: given %d", maxMsg)
	}

	var waitTimeSeconds int32 = 10

	// Must always be above `WaitTimeSeconds` otherwise `ReceiveMessageWithContext`
	// trigger context timeout error.
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(waitTimeSeconds+5))
	defer cancel()

	res, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueURL),
		MaxNumberOfMessages:   10,
		WaitTimeSeconds:       waitTimeSeconds,
		MessageAttributeNames: []string{"All"},
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	return res.Messages, nil
}

func (s Queuer) Delete(ctx context.Context, queueURL, rcvHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(rcvHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
