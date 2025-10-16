package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
			Brokers:  []string{broker},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) ReadMessage(ctx context.Context) (*JobMessage, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
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
