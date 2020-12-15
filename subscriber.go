package htsqs

import (
	"errors"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jpillora/backoff"
)

const (
	// defaultMaxMessagesPerBatch is default the number of messages
	// the subscriber will attempt to fetch on each receive.
	defaultMaxMessagesPerBatch int64 = 10

	// defaultWaitTimeoutSeconds the duration (in seconds) for which the call waits for a message to arrive
	// in the queue before returning. If a message is available, the call returns
	// sooner than TimeSeconds. If no messages are available and the wait time
	// expires, the call returns successfully with an empty list of messages.
	defaultWaitTimeoutSeconds int64 = 10

	// defaultVisibilityTimeout The duration (in seconds) that the received messages are hidden from subsequent
	// retrieve requests after being retrieved by a ReceiveMessage request.
	defaultVisibilityTimeout int64 = 30

	// defaultNumConsumers is the number of consumers per subscriber
	defaultNumConsumers int = 3
)

type atomicBool int32

func (b *atomicBool) isSet() bool {
	return atomic.LoadInt32((*int32)(b)) != 0
}

func (b *atomicBool) setTrue() error {
	if atomic.CompareAndSwapInt32((*int32)(b), 0, 1) {
		return nil
	}
	return errors.New("value is already set")
}

// receiver is the interface to sqsiface.SQSAPI. The only purpose is to be able to mock sqs for testing. See mock_test.go
type receiver interface {
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
	ChangeMessageVisibility(params *sqs.ChangeMessageVisibilityInput) (*sqs.ChangeMessageVisibilityOutput, error)
}

// Logger interface allows to use other loggers than standard log.Logger
type Logger interface {
	Printf(string, ...interface{})
}

// SubscriberConfig holds the info required to work with Amazon SQS
type SubscriberConfig struct {

	// AWS session
	AWSSession *session.Session

	// SQS queue from which the subscriber is going to consume from
	SqsQueueURL string

	// number of messages the subscriber will attempt to fetch on each receive.
	MaxMessagesPerBatch int64

	// the duration (in seconds) for which the call waits for a message to arrive
	// in the queue before returning. If a message is available, the call returns
	// sooner than TimeSeconds. If no messages are available and the wait time
	// expires, the call returns successfully with an empty list of messages.
	TimeoutSeconds int64

	// The duration (in seconds) that the received messages are hidden from subsequent
	// retrieve requests after being retrieved by a ReceiveMessage request.
	// VisibilityTimeout should be < time needed to process a message
	VisibilityTimeout int64

	// number of consumers per subscriber
	NumConsumers int

	// subscriber logger
	Logger Logger
}

// Subscriber is an SQS client that allows a user to
// consume messages from AWS SQS.
// Once Stop has been called on subscriber, it might not be reused;
// future calls to methods such as Consume or Stop will return an error.
type Subscriber struct {
	sqs      receiver
	cfg      SubscriberConfig
	stopped  atomicBool
	consumed atomicBool
	stop     chan error
}

// Consume starts consuming messages from the SQS queue.
// Returns a channel of SubscriberMessage to consume them and a channel of errors
func (s *Subscriber) Consume() (<-chan *SQSMessage, <-chan error, error) {
	if s.stopped.isSet() {
		return nil, nil, errors.New("SQS subscriber is already stopped")
	}

	if s.consumed.setTrue() != nil {
		return nil, nil, errors.New("SQS subscriber is already running")
	}

	var wg sync.WaitGroup
	var messages chan *SQSMessage
	var errCh chan error

	messages = make(chan *SQSMessage, s.cfg.MaxMessagesPerBatch*int64(s.cfg.NumConsumers))
	errCh = make(chan error, int64(s.cfg.NumConsumers))

	backoffCounter := backoff.Backoff{
		Factor: 1,
		Min:    time.Second,
		Max:    30 * time.Second,
		Jitter: true,
	}

	for i := 1; i <= s.cfg.NumConsumers; i++ {
		wg.Add(1)
		go func(workerID int, backoffCfg backoff.Backoff) {
			s.cfg.Logger.Printf("Consumer %d listening for messages", workerID)
			defer wg.Done()

			var msgs *sqs.ReceiveMessageOutput
			var err error

			for !s.stopped.isSet() {
				msgs, err = s.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
					MessageAttributeNames: []*string{aws.String(sqs.QueueAttributeNameAll)},
					MaxNumberOfMessages:   &s.cfg.MaxMessagesPerBatch,
					QueueUrl:              &s.cfg.SqsQueueURL,
					WaitTimeSeconds:       &s.cfg.TimeoutSeconds,
					VisibilityTimeout:     &s.cfg.VisibilityTimeout,
				})

				if err != nil {
					// Error found, send the error
					errCh <- err
					time.Sleep(backoffCfg.Duration())
					continue
				}

				s.cfg.Logger.Printf("Found %d messages\n", len(msgs.Messages))
				backoffCfg.Reset()
				// for each message, pass to output
				for _, msg := range msgs.Messages {
					messages <- &SQSMessage{
						sub:        s,
						RawMessage: msg,
					}
				}
			}
		}(i, backoffCounter)
	}

	go func() {
		wg.Wait()
		close(messages)
		close(errCh)
		s.stop <- nil
		close(s.stop)
	}()

	s.cfg.Logger.Printf("SQS subscriber listening for messages\n")
	return messages, errCh, nil
}

// Stop stop gracefully the Subscriber
func (s *Subscriber) Stop() error {
	if err := s.stopped.setTrue(); err != nil {
		return errors.New("SQS subscriber is already stopped")
	}
	return <-s.stop
}

func defaultSubscriberConfig(cfg *SubscriberConfig) {
	if cfg.AWSSession == nil {
		cfg.AWSSession = session.Must(session.NewSession())
	}

	if cfg.MaxMessagesPerBatch == 0 {
		cfg.MaxMessagesPerBatch = defaultMaxMessagesPerBatch
	}

	if cfg.TimeoutSeconds == 0 {
		cfg.TimeoutSeconds = defaultWaitTimeoutSeconds
	}

	if cfg.VisibilityTimeout == 0 {
		cfg.VisibilityTimeout = defaultVisibilityTimeout
	}

	if cfg.NumConsumers == 0 {
		cfg.NumConsumers = defaultNumConsumers
	}

	if cfg.Logger == nil {
		cfg.Logger = log.New(os.Stdout, "", log.LstdFlags|log.LUTC)
	}
}

// NewSubscriber creates a new AWS SQS subscriber
func NewSubscriber(cfg SubscriberConfig) *Subscriber {
	defaultSubscriberConfig(&cfg)
	return &Subscriber{cfg: cfg, sqs: sqs.New(cfg.AWSSession), stop: make(chan error, 1)}
}
