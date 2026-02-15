.PHONY: run build test lint

run:
	go run ./cmd/server

build:
	go build -o url-shortener ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run
