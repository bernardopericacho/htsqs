package htsqs

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/require"
)

func TestSubscriber(t *testing.T) {
	numMessages := 10
	queue := make(chan *SQSMessage)
	defer close(queue)
	subs := NewSubscriber(SubscriberConfig{})
	subs.sqs = &sqsMock{queue: queue}

	stopErrChannel := make(chan error)

	go func() {
		for i := 0; i < numMessages; i++ {
			message := fmt.Sprintf("Message: %d", i)
			queue <- &SQSMessage{subs, &sqs.Message{Body: &message}}
		}
		stopErrChannel <- subs.Stop()
		close(stopErrChannel)
	}()

	messages, errch, err := subs.Consume()
	require.NoError(t, err)

	// try to start consuming again while we are still consuming
	_, _, err = subs.Consume()
	require.EqualError(t, err, "SQS subscriber is already running")

	i := 0
	for m := range messages {
		require.True(t, strings.HasPrefix(string(m.Message()), "Message: "))
		require.NoError(t, m.ChangeMessageVisibility(aws.Int64(43200)))
		require.NoError(t, m.Done())
		i++
	}
	require.NoError(t, <-stopErrChannel)
	require.NoError(t, <-errch)
	require.Equal(t, numMessages, i)
	require.EqualError(t, subs.Stop(), "SQS subscriber is already stopped")

	// try to start consuming again when the consumer has already been used
	_, _, err = subs.Consume()
	require.EqualError(t, err, "SQS subscriber is already stopped")
}

func TestSubscriberDefaults(t *testing.T) {

	tt := []struct {
		name                  string
		sqsConfig             SubscriberConfig
		expectedAfterDefaults SubscriberConfig
	}{
		{
			"Custom parameters",
			SubscriberConfig{AWSSession: session.Must(session.NewSession()), MaxMessagesPerBatch: 1, TimeoutSeconds: 1, VisibilityTimeout: 1, NumConsumers: 1, Logger: log.New(os.Stderr, "", log.LstdFlags)},
			SubscriberConfig{MaxMessagesPerBatch: 1, TimeoutSeconds: 1, VisibilityTimeout: 1, NumConsumers: 1, Logger: log.New(os.Stderr, "", log.LstdFlags)},
		},
		{
			"Use defaults parameters",
			SubscriberConfig{},
			SubscriberConfig{MaxMessagesPerBatch: 10, TimeoutSeconds: 10, VisibilityTimeout: 30, NumConsumers: 3, Logger: log.New(os.Stdout, "", log.LstdFlags|log.LUTC)},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Check provided config is not modified
			NewSubscriber(tc.sqsConfig)
			require.Exactly(t, tc.sqsConfig, tc.sqsConfig)

			// Check if defaults are properly calculated
			initialAWSSession := tc.sqsConfig.AWSSession
			defaultSubscriberConfig(&tc.sqsConfig)
			// Check AWS session conf
			if initialAWSSession == nil {
				require.NotNil(t, tc.sqsConfig.AWSSession)
				tc.sqsConfig.AWSSession = nil
			} else {
				require.Equal(t, initialAWSSession, tc.sqsConfig.AWSSession)
				tc.expectedAfterDefaults.AWSSession = initialAWSSession
			}
			require.Exactly(t, tc.sqsConfig, tc.expectedAfterDefaults)

		})
	}
}
