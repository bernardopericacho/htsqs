package htsqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessage is the SQS implementation of `SubscriberMessage`.
type SQSMessage struct {
	sub        *Subscriber
	RawMessage *sqs.Message
}

// Message returns the body of the SQS message in bytes
func (m *SQSMessage) Message() []byte {
	return []byte(*m.RawMessage.Body)
}

// Done deletes the message from SQS.
func (m *SQSMessage) Done() error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      &m.sub.cfg.SqsQueueURL,
		ReceiptHandle: m.RawMessage.ReceiptHandle,
	}
	_, err := m.sub.sqs.DeleteMessage(deleteParams)
	return err
}
