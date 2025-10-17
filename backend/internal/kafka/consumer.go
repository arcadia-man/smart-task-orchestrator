package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

type JobMessage struct {
	JobID   string         `json:"jobId"`
	Payload map[string]any `json:"payload"`
}

func NewConsumer(broker, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     []string{broker},
			Topic:       topic,
			GroupID:     groupID,
			MinBytes:    1,                 // 1 byte minimum
			MaxBytes:    10e6,              // 10MB
			StartOffset: kafka.FirstOffset, // Start from beginning for new consumer groups
		}),
	}
}

func (c *Consumer) ReadMessage(ctx context.Context) (*JobMessage, error) {
	// Add timeout to prevent indefinite blocking
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	msg, err := c.reader.ReadMessage(ctxWithTimeout)
	if err != nil {
		// Check if it's a timeout error
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("timeout waiting for message")
		}
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	var jobMsg JobMessage
	if err := json.Unmarshal(msg.Value, &jobMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	log.Printf("📨 Received job message: %s", jobMsg.JobID)
	return &jobMsg, nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
