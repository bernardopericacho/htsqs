package htsqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type sqsMock struct {
	queue <-chan *SQSMessage
}

func (s *sqsMock) ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	select {
	case message := <-s.queue:
		if message == nil {
			s.queue = nil
			return &sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}, nil
		}
		stringMessage := string(message.Message())
		return &sqs.ReceiveMessageOutput{Messages: []*sqs.Message{{Body: &stringMessage, ReceiptHandle: &stringMessage}}}, nil
	default:
		return &sqs.ReceiveMessageOutput{Messages: []*sqs.Message{}}, nil
	}
}

func (s *sqsMock) DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return nil, nil
}

func (s *sqsMock) ChangeMessageVisibility(*sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error) {
	return nil, nil
}
