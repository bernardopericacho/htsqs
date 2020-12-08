# HTSQS

[![PkgGoDev](https://pkg.go.dev/badge/golang.org/x/tools)](https://pkg.go.dev/github.com/bernardopericacho/htsqs) [![Build Status](https://travis-ci.com/bernardopericacho/htsqs.svg?branch=master)](https://travis-ci.com/bernardopericacho/htsqs)

HTSQS is a high throughput golang AWS SQS consumer.

## Install

`go get -u github.com/bernardopericacho/htsqs`

## Features

* **High throughput** - ability multiple consumers that concurrently receive messages from AWS SQS and push them into a single channel for consumption
* **Late ACK** - mechanism for acknowledging messages once they've been processed
* **Error processing** - error processing to decide whether to stop consuming and exponential backoff setup when errors occur
* **Graceful shutdown**

## License
This project is licensed under [MIT License](./LICENSE).

