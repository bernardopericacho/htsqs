package transport

// Subscriber is the interface clients can use to receive messages from SQS
type Subscriber interface {
	// Consume will return a channel of raw messages to be consumed and a channel of errors
	Consume() (<-chan SubscriberMessage, <-chan error, error)
	// Stop will initiate a graceful shutdown of the subscriber connection.
	Stop() error
}

// SubscriberMessage is an interface to encapsulate subscriber messages and provide
// a mechanism for acknowledging messages once they've been processed.
type SubscriberMessage interface {
	Message() []byte
	Done() error
}
