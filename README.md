# dog-identifier-api

Locally hosted API written in Go that determines whether a valid image is of a dog or not via Google Vision API.

The purpose of the project is for the author to learn about Go, as well as some simple usage of a Google Cloud API.

The project only supports image identification via URL's, but is built so functionality can be added where needed.

This project also includes HTTPS-enabled hosting (given proper certificate files from environment variables) and testing.

## Project setup

```bash
go get .
```

## Local deployment

```bash
go run main.go
```

## Run all handler tests

```bash
go test ./handlers
```

## Run specific test

```bash
go test *.go *_test.go
```
