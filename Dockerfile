# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /build/server ./server
COPY migrations/ ./migrations/

EXPOSE 8080

CMD ["./server"]
