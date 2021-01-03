package sns

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// sender is the interface to sns.SNS. Its sole purpose is to make
// Publisher.service and interface that we can mock for testing.
type sender interface {
	PublishWithContext(ctx context.Context, input *sns.PublishInput, o ...request.Option) (*sns.PublishOutput, error)
}

// Config holds the info required to work with AWS SNS
type Config struct {

	// AWS session
	AWSSession *session.Session

	// Topic ARN where the messages are going to be sent
	TopicArn string
}

// Publisher is the AWS SNS message publisher
type Publisher struct {
	sns sender
	cfg Config
}

// Publish allows SNS Publisher to implement the publisher.Publisher interface
// and publish messages to an AWS SNS backend
func (p *Publisher) Publish(ctx context.Context, msg json.Marshaler) error {
	b, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(b)),
		TopicArn: &p.cfg.TopicArn,
	}

	if err := input.Validate(); err != nil {
		return err
	}
	_, err = p.sns.PublishWithContext(ctx, input)

	return err
}

func defaultPublisherConfig(cfg *Config) {
	if cfg.AWSSession == nil {
		cfg.AWSSession = session.Must(session.NewSession())
	}
}

// New creates a new AWS SNS publisher
func New(cfg Config) *Publisher {
	defaultPublisherConfig(&cfg)
	return &Publisher{cfg: cfg, sns: sns.New(cfg.AWSSession)}
}
