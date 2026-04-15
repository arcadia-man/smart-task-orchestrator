package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(ctx context.Context, key string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: payload,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func(kafka.Message) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if err := handler(m); err != nil {
			// In a real app, we might want to log this or send to DLQ
			continue
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
