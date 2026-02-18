# Build stage
FROM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /build/server ./server
COPY migrations/ ./migrations/
COPY static/ ./static/

EXPOSE 8080

CMD ["./server"]
