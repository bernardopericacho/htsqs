# HTSQS

[![Latest Version](http://img.shields.io/github/v/release/bernardopericacho/htsqs.svg)](https://github.com/bernardopericacho/htsqs/releases) [![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/tools)](https://pkg.go.dev/github.com/bernardopericacho/htsqs) [![Build Status](https://travis-ci.com/bernardopericacho/htsqs.svg?branch=master)](https://travis-ci.com/bernardopericacho/htsqs)

HTSQS is a high throughput golang AWS SQS consumer.

## Install

`go get -u github.com/bernardopericacho/htsqs`

## Features

* **High throughput** - ability to set multiple consumers that concurrently receive messages from AWS SQS and push them into a single channel for consumption
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
    
    "github.com/bernardopericacho/htsqs"
)

func main() {
    // Create a new subscriber, assuming we are configuring our credentials through environment variables
    subs := htsqs.NewSubscriber(htsqs.SubscriberConfig{SqsQueueURL: <MY_SQS_QUEUE_URL>})
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

## License

This project is licensed under [MIT License](./LICENSE).

