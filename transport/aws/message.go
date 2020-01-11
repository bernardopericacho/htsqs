package aws

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// sqsMessage is the SQS implementation of `SubscriberMessage`.
type sqsMessage struct {
	sub     *Subscriber
	message *sqs.Message
}

func (m *sqsMessage) Message() []byte {
	return []byte(*m.message.Body)
}

func (m *sqsMessage) Done() error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      &m.sub.cfg.SqsQueueUrl,
		ReceiptHandle: m.message.ReceiptHandle,
	}
	_, err := m.sub.sqs.DeleteMessage(deleteParams)
	return err
}
