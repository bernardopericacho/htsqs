package sns

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/require"
)

type jsonString string

func (js jsonString) MarshalJSON() ([]byte, error) {
	return []byte(js), nil
}

func TestPublisher(t *testing.T) {
	queue := make(chan *string, 1)
	defer close(queue)
	pubs := New(Config{})
	pubs.sqs = &sqsPublisherMock{queue: queue}

	testString := jsonString(`{"msg":"message"}`)
	require.NoError(t, pubs.Publish(context.Background(), testString))
	publishedMessage := <-queue
	require.Equal(t, *publishedMessage, `{"msg":"message"}`)
}

func TestPublisherDefaults(t *testing.T) {

	tt := []struct {
		name                  string
		sqsConfig             Config
		expectedAfterDefaults Config
	}{
		{
			"Custom parameters",
			Config{AWSSession: session.Must(session.NewSession()), QueueURL: "myQueueURL"},
			Config{QueueURL: "myQueueURL"},
		},
		{
			"Use defaults parameters",
			Config{},
			Config{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Check provided config is not modified
			New(tc.sqsConfig)
			require.Exactly(t, tc.sqsConfig, tc.sqsConfig)

			// Check if defaults are properly calculated
			initialAWSSession := tc.sqsConfig.AWSSession
			defaultPublisherConfig(&tc.sqsConfig)
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
