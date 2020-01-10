package aws

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/require"

	"htsqs/transport"
)

func TestSubscriber(t *testing.T) {
	numMessages := 10
	queue := make(chan transport.SubscriberMessage)
	defer close(queue)
	subs := NewSubscriber(session.Must(session.NewSession()), SubscriberConfig{})
	subs.sqs = &sqsMock{queue: queue}

	stopErrChannel := make(chan error)

	go func() {
		for i := 0; i < numMessages; i++ {
			message := fmt.Sprintf("Message: %d", i)
			queue <- &sqsMessage{subs, &sqs.Message{Body: &message}}
		}
		stopErrChannel <- subs.Stop()
		close(stopErrChannel)
	}()

	i := 0
	messages, errch, err := subs.Consume()
	require.NoError(t, err)

	// try to start consuming again while we are still consuming
	_, _, err = subs.Consume()
	require.EqualError(t, err, "sqs subscriber is already running")

	for m := range messages {
		require.Equal(t, fmt.Sprintf("Message: %d", i), string(m.Message()))
		require.NoError(t, m.Done())
		i++
	}
	require.NoError(t, <-stopErrChannel)
	require.NoError(t, <-errch)
	require.Equal(t, numMessages, i)
	require.EqualError(t, subs.Stop(), "sqs subscriber is already stopped")

	// try to start consuming again when the consumer has already been used
	_, _, err = subs.Consume()
	require.EqualError(t, err, "sqs subscriber is already stopped")
}

func TestSubscriberDefaults(t *testing.T) {

	tt := []struct {
		name                  string
		sqsConfig             SubscriberConfig
		expectedAfterDefaults SubscriberConfig
	}{
		{
			"Custom parameters",
			SubscriberConfig{MaxMessagesPerBatch: 1, TimeoutSeconds: 1, VisibilityTimeout: 1, Logger: log.New(os.Stderr, "", log.LstdFlags)},
			SubscriberConfig{MaxMessagesPerBatch: 1, TimeoutSeconds: 1, VisibilityTimeout: 1, Logger: log.New(os.Stderr, "", log.LstdFlags)},
		},
		{
			"Use defaults parameters",
			SubscriberConfig{},
			SubscriberConfig{MaxMessagesPerBatch: 10, TimeoutSeconds: 10, VisibilityTimeout: 30, Logger: log.New(os.Stdout, "", log.LstdFlags|log.LUTC)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Check provided config is not modified
			NewSubscriber(session.Must(session.NewSession()), tc.sqsConfig)
			require.Exactly(t, tc.sqsConfig, tc.sqsConfig)

			// Check if defaults are properly calculated
			defaultSubscriberConfig(&tc.sqsConfig)
			require.Exactly(t, tc.sqsConfig, tc.expectedAfterDefaults)

		})
	}
}
