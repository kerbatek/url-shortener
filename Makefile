.PHONY: run build test lint docker-up docker-down

run:
	go run ./cmd/server

build:
	go build -o url-shortener ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run

docker-up:
	docker compose up --build

docker-down:
	docker compose down
