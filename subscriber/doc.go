// Package subscriber provides the functionalities to consume messages from an AWS SQS queue.
// For more information about to AWS SQS go to https://aws.amazon.com/sqs/
//
// AWS SQS Subscriber
//
// Subscriber is a high throughput golang AWS SQS client that can create multiple consumers
// that concurrently receive messages from AWS SQS and push them into a single channel for consumption.
//
// Worker
//
// Worker is the service implementation of a Subscriber.
package subscriber
