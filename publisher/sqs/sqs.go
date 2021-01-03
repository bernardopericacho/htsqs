package sns

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// sender is the interface to sqs.SQS. Its sole purpose is to make
// Publisher.service and interface that we can mock for testing.
type sender interface {
	SendMessageWithContext(ctx aws.Context, input *sqs.SendMessageInput, opts ...request.Option) (*sqs.SendMessageOutput, error)
}

// Config holds the info required to work with AWS SQS to publish a message
type Config struct {

	// AWS session
	AWSSession *session.Session

	// SQS queue where the publisher is going to push messages to
	QueueURL string
}

// Publisher is the AWS SNS message publisher
type Publisher struct {
	sqs sender
	cfg Config
}

// Publish allows SQS Publisher to implement the publisher.Publisher interface
// and publish messages to an AWS SQS backend
func (p *Publisher) Publish(ctx context.Context, msg json.Marshaler) error {
	b, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(string(b)),
		QueueUrl:    &p.cfg.QueueURL,
	}

	if err := input.Validate(); err != nil {
		return err
	}
	_, err = p.sqs.SendMessageWithContext(ctx, input)

	return err
}

func defaultPublisherConfig(cfg *Config) {
	if cfg.AWSSession == nil {
		cfg.AWSSession = session.Must(session.NewSession())
	}
}

// New creates a new AWS SQS publisher
func New(cfg Config) *Publisher {
	defaultPublisherConfig(&cfg)
	return &Publisher{cfg: cfg, sqs: sqs.New(cfg.AWSSession)}
}
