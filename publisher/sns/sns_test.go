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
	pubs.sns = &snsPublisherMock{queue: queue}

	testString := jsonString(`{"msg":"message"}`)
	require.NoError(t, pubs.Publish(context.Background(), testString))
	publishedMessage := <-queue
	require.Equal(t, *publishedMessage, `{"msg":"message"}`)
}

func TestPublisherDefaults(t *testing.T) {

	tt := []struct {
		name                  string
		snsConfig             Config
		expectedAfterDefaults Config
	}{
		{
			"Custom parameters",
			Config{AWSSession: session.Must(session.NewSession()), TopicArn: "myTopicARN"},
			Config{TopicArn: "myTopicARN"},
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
			New(tc.snsConfig)
			require.Exactly(t, tc.snsConfig, tc.snsConfig)

			// Check if defaults are properly calculated
			initialAWSSession := tc.snsConfig.AWSSession
			defaultPublisherConfig(&tc.snsConfig)
			// Check AWS session conf
			if initialAWSSession == nil {
				require.NotNil(t, tc.snsConfig.AWSSession)
				tc.snsConfig.AWSSession = nil
			} else {
				require.Equal(t, initialAWSSession, tc.snsConfig.AWSSession)
				tc.expectedAfterDefaults.AWSSession = initialAWSSession
			}
			require.Exactly(t, tc.snsConfig, tc.expectedAfterDefaults)

		})
	}
}
