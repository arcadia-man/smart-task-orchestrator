FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api/main.go
RUN go build -o worker ./cmd/worker/main.go
RUN go build -o scheduler ./cmd/scheduler/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/worker .
COPY --from=builder /app/scheduler .

# Install docker client in alpine to talk to the socket
RUN apk add --no-cache docker-cli

EXPOSE 8080

CMD ["./api"]