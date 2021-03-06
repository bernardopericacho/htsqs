package subscriber

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSMessage is the implementation of a SQS message
type SQSMessage struct {
	sub        *Subscriber
	rawMessage *sqs.Message
}

// Body returns the body of the SQS message in bytes
func (m *SQSMessage) Body() []byte {
	return []byte(*m.rawMessage.Body)
}

// MessageAttributes returns the message attributes
func (m *SQSMessage) MessageAttributes() map[string]*sqs.MessageAttributeValue {
	return m.rawMessage.MessageAttributes
}

// Done deletes the message from SQS.
func (m *SQSMessage) Done() error {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      &m.sub.cfg.SqsQueueURL,
		ReceiptHandle: m.rawMessage.ReceiptHandle,
	}
	_, err := m.sub.sqs.DeleteMessage(deleteParams)
	return err
}

// ChangeMessageVisibility modifies current message visibility timeout to the one specified in the parameters.
// This is normally useful when the message processing is taking more time than the default visibility timeout
func (m *SQSMessage) ChangeMessageVisibility(newVisibilityTimeout *int64) error {
	changeVisibilityParams := &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &m.sub.cfg.SqsQueueURL,
		ReceiptHandle:     m.rawMessage.ReceiptHandle,
		VisibilityTimeout: newVisibilityTimeout,
	}

	// Validate values
	if err := changeVisibilityParams.Validate(); err != nil {
		return err
	}

	_, err := m.sub.sqs.ChangeMessageVisibility(changeVisibilityParams)
	return err
}
