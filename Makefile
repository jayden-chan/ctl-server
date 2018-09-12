default: lint test

build:
	go fmt ./...
	go build -race

lint:
	go fmt ./...
	golint ./...
	go vet ./...

.PHONY: build lint