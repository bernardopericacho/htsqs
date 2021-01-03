# HTSQS

[![Latest Version](http://img.shields.io/github/v/release/bernardopericacho/htsqs.svg)](https://github.com/bernardopericacho/htsqs/releases) [![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/tools)](https://pkg.go.dev/github.com/bernardopericacho/htsqs) [![Build Status](https://travis-ci.com/bernardopericacho/htsqs.svg?branch=master)](https://travis-ci.com/bernardopericacho/htsqs)

HTSQS is a high throughput golang AWS SQS consumer.

## Install

`go get -u github.com/bernardopericacho/htsqs`

## Features

* **High throughput** - a subscriber has the ability to create multiple consumers that concurrently receive messages from AWS SQS and push them into a single channel for consumption
* **Late ACK** - mechanism for acknowledging messages once they have been processed
* **Message visibility** modify message visibility
* **Error processing** - error processing to decide whether to stop consuming and exponential backoff setup when errors occur
* **Graceful shutdown**

## Getting started

### Consume messages from an AWS SQS Queue 

```go
package main

import (
    
    "log"
    
    "github.com/bernardopericacho/htsqs/subscriber"
)

func main() {
    // Create a new subscriber, assuming we are configuring our credentials following 
	// environment variables or IAM Roles: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html
    subs := subscriber.New(subscriber.Config{SqsQueueURL: <MY_SQS_QUEUE_URL>})
    // Call consume
    messagesCh, errCh, err := subs.Consume()
    if err != nil {
        log.Fatal("Error when trying to consume messages from the SQS Queue")
    }
    // Loop over all the messages or errors
    for {
        select {
            case msg := <-messagesCh:
                log.Println("received new message", msg)
            case err := <-errCh:
                log.Println("received new error", err)
        }
    }
}
```

### Create a worker service to consume from a SQS queue

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/bernardopericacho/htsqs/subscriber"
)

func main() {
    cfg := subscriber.WorkerConfig{
		Subscriber: subscriber.New(subscriber.Config{
			SqsQueueURL: "",
		}),
	}
	
	worker := subscriber.NewWorker(cfg)
	ctx := context.Background()
	if err := worker.Start(ctx); err != subscriber.ErrWorkerClosed {
		stopErr := worker.Stop()
		if stopErr != nil {
			log.Printf("Worker start failed: %v\n", fmt.Errorf("%s: %w", stopErr.Error(), err))
		} else {
			log.Println("Worker start failed")
		}
	}
}

```

## License

This project is licensed under [MIT License](./LICENSE).

