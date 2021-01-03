package subscriber

import (
	"context"
	"errors"
	"log"
)

// ErrWorkerClosed is returned by the Worker 'Start' method after a call to 'Stop'.
var ErrWorkerClosed = errors.New("worker closed")

func defaultMessageHandler(ctx context.Context, w *Worker, m *SQSMessage) {
	log.Printf("Message received: '%s'", string(m.Message()))
	if err := m.Done(); err != nil {
		log.Printf("Error when deleting message from SQS: %v", err)
	}
}

func defaultErrorHandler(ctx context.Context, w *Worker, e error) {
	log.Printf("Error when receiving messages from SQS: %v", e)
	w.lastErr <- e
}

// WorkerConfig is the worker startup config
type WorkerConfig struct {

	// SQS subscriber
	Subscriber *Subscriber

	// SQS Message Handler
	MessageHandler func(context.Context, *Worker, *SQSMessage)

	// SQS Error Handler
	ErrorHandler func(context.Context, *Worker, error)
}

func defaultWorkerConfig(cfg *WorkerConfig) {
	if cfg.MessageHandler == nil {
		cfg.MessageHandler = defaultMessageHandler
	}
	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = defaultErrorHandler
	}
}

// Worker represents a SQS worker service
type Worker struct {
	lastErr chan error
	config  *WorkerConfig
}

// Start triggers the process to start consuming messages from the SQS subscriber.
// Blocks until `lastErr` is set or `Stop()` is called
func (w *Worker) Start(ctx context.Context) error {

	sqsMessages, errorCh, err := w.config.Subscriber.Consume()
	if err != nil {
		return err
	}

	// Process all errors in a goroutine
	go func() {
		for err := range errorCh {
			w.config.ErrorHandler(ctx, w, err)
		}
	}()

	// Process each message in a goroutine
	for message := range sqsMessages {
		go w.config.MessageHandler(ctx, w, message)
	}

	return <-w.lastErr
}

// Stop gracefully stops the subscriber.
func (w *Worker) Stop() error {
	if err := w.config.Subscriber.Stop(); err != nil {
		return err
	}
	w.lastErr <- ErrWorkerClosed
	close(w.lastErr)
	return nil
}

// Config returns current configuration
func (w *Worker) Config() *WorkerConfig {
	return w.config
}

// NewWorker creates a new Worker based on the given configuration that process messages from AWS SQS
func NewWorker(conf WorkerConfig) *Worker {
	defaultWorkerConfig(&conf)
	return &Worker{lastErr: make(chan error, 1), config: &conf}
}
