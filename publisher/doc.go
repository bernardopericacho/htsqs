// Package publisher provides the functionalities to publish messages to AWS SNS topic or to a AWS SQS queue that
// will be consumed by an SQS subscriber. https://github.com/bernardopericacho/htsqs/subscriber
//
// AWS SNS Publisher
//
// Publish messages to the given AWS SNS Topic. Using AWS SNS + SQS (publish/subscribe solution),
// AWS SNS publisher will publish messages to the SNS topic and they will be broadcast to all the connected SQS queues.
// For more information about to AWS SQS go to https://aws.amazon.com/sns/
//
// AWS SQS Publisher
//
// Publish messages to the given AWS SQS Queue. AWS SQS publisher will publish messages to the AWS SQS queue
// for asynchronous message processing.
// For more information about to AWS SQS go to https://aws.amazon.com/sqs/
package publisher
