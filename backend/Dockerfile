# Multi-stage build for Go backend
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scheduler ./cmd/scheduler  
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o worker ./cmd/worker

# Final stage
FROM alpine:latest

# Install ca-certificates and docker client
RUN apk --no-cache add ca-certificates docker-cli

WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/api .
COPY --from=builder /app/scheduler .
COPY --from=builder /app/worker .

# Default command (can be overridden in docker-compose)
CMD ["./api"]