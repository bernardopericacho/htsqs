package sns

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsPublisherMock struct {
	queue chan<- *string
}

func (p *sqsPublisherMock) SendMessageWithContext(ctx context.Context, input *sqs.SendMessageInput, opts ...request.Option) (*sqs.SendMessageOutput, error) {
	p.queue <- input.MessageBody
	return &sqs.SendMessageOutput{}, nil
}
