.PHONY: run build test lint docker-dev-up docker-prod-up docker-down

run:
	go run ./cmd/server

build:
	go build -o url-shortener ./cmd/server

test:
	go test ./...

lint:
	golangci-lint run

docker-dev-up:
	docker compose up --build

docker-prod-up:
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up

docker-down:
	docker compose down
