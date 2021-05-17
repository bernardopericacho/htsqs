package subscriber

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/require"
)

func TestWorker(t *testing.T) {
	c := new(WorkerConfig)

	queue := make(chan *SQSMessage)
	defer close(queue)
	subs := New(Config{})
	subs.sqs = &sqsMock{queue: queue}
	c.Subscriber = subs
	worker := NewWorker(*c)

	errsChannelStop := make(chan error)

	// Send messages to the message channel to be consumed
	go func() {
		for i := 0; i < 10; i++ {
			message := fmt.Sprintf("Message: %d", i)

			queue <- &SQSMessage{
				sub: subs,
				rawMessage: &sqs.Message{
					Body: &message,
				},
			}
		}
		errsChannelStop <- worker.Stop()
		close(errsChannelStop)
	}()

	require.Equal(t, ErrWorkerClosed, worker.Start(context.TODO()))
	require.NoError(t, <-errsChannelStop)
	require.EqualError(t, worker.Stop(), "SQS subscriber is already stopped")
	require.EqualError(t, worker.Start(context.TODO()), "SQS subscriber is already stopped")
}

func TestWorkerAlreadyRunning(t *testing.T) {
	c := new(WorkerConfig)

	queue := make(chan *SQSMessage)
	defer close(queue)
	subs := New(Config{})
	subs.sqs = &sqsMock{queue: queue}
	c.Subscriber = subs
	worker := NewWorker(*c)

	errsChannelStart := make(chan error)
	errsChannelStop := make(chan error)

	go func() {
		message := fmt.Sprintf("Message")

		queue <- &SQSMessage{
			sub: subs,
			rawMessage: &sqs.Message{
				Body: &message,
			},
		}
		errsChannelStop <- worker.Stop()
		close(errsChannelStop)
	}()

	go func() {
		errsChannelStart <- worker.Start(context.TODO())
		close(errsChannelStart)
	}()

	require.Equal(t, ErrWorkerClosed, worker.Start(context.TODO()))
	require.EqualError(t, <-errsChannelStart, "SQS subscriber is already running")
	require.NoError(t, <-errsChannelStop)
}

func TestWorkerError(t *testing.T) {

	c := new(WorkerConfig)
	errorQueue := make(chan error)
	defer close(errorQueue)

	subs := New(Config{})
	subs.sqs = &sqsMock{errorQueue: errorQueue}
	c.Subscriber = subs
	worker := NewWorker(*c)

	errsChannelStart := make(chan error)
	go func() {
		errsChannelStart <- worker.Start(context.TODO())
		close(errsChannelStart)
	}()

	AWSError := errors.New("AWS very bad error")
	errorQueue <- AWSError
	require.NoError(t, worker.Stop())
	require.EqualError(t, <-errsChannelStart, AWSError.Error())

}
