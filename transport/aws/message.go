package aws

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessage is the SQS implementation of `SubscriberMessage`.
type SQSMessage struct {
	sub        *Subscriber
	RawMessage *sqs.Message
}

func (m *SQSMessage) Message() []byte {
	return []byte(*m.RawMessage.Body)
}

func (m *SQSMessage) Done() error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      &m.sub.cfg.SqsQueueUrl,
		ReceiptHandle: m.RawMessage.ReceiptHandle,
	}
	_, err := m.sub.sqs.DeleteMessage(deleteParams)
	return err
}
