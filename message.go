package htsqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessage is the implementation of a SQS message
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

// ChangeMessageVisibility modifies current message visibility timeout to the one specified in the parameters.
// This is normally useful when the message processing is taking more time than the default visibility timeout
func (m *SQSMessage) ChangeMessageVisibility(newVisibilityTimeout *int64) error {
	changeVisibilityParams := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &m.sub.cfg.SqsQueueURL,
		ReceiptHandle:     m.RawMessage.ReceiptHandle,
		VisibilityTimeout: newVisibilityTimeout,
	}

	// Validate values
	if err := changeVisibilityParams.Validate(); err != nil {
		return err
	}

	_, err := m.sub.sqs.ChangeMessageVisibility(changeVisibilityParams)
	return err
}
