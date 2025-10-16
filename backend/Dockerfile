FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o api cmd/api/main.go
RUN go build -o worker cmd/worker/main.go
RUN go build -o scheduler cmd/scheduler/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/api .
COPY --from=builder /app/worker .
COPY --from=builder /app/scheduler .

# Make binaries executable
RUN chmod +x api worker scheduler

# Default command (can be overridden)
CMD ["./api"]