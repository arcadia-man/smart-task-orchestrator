package kafka

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/segmentio/kafka-go"
)

type Producer struct {
    writer *kafka.Writer
}

type JobExecutionMessage struct {
    RunID       string                 `json:"run_id"`
    SchedulerID string                 `json:"scheduler_id"`
    Generation  int                    `json:"generation"`
    Image       string                 `json:"image"`
    Command     string                 `json:"command"`
    Env         map[string]string      `json:"env"`
    TriggeredBy string                 `json:"triggered_by"`
    RunAt       time.Time              `json:"run_at"`
    Metadata    map[string]interface{} `json:"metadata"`
}

func NewProducer(brokers string, topic string) *Producer {
    writer := &kafka.Writer{
        Addr:         kafka.TCP(brokers),
        Topic:        topic,
        Balancer:     &kafka.Hash{}, // Use hash balancer for consistent partitioning
        RequiredAcks: kafka.RequireOne,
        Async:        false, // Synchronous for reliability
        BatchTimeout: 10 * time.Millisecond,
        BatchSize:    100,
    }

    log.Printf("✅ Kafka producer initialized for topic: %s", topic)
    return &Producer{writer: writer}
}

func (p *Producer) PublishJobExecution(ctx context.Context, msg JobExecutionMessage) error {
    // Serialize message
    data, err := json.Marshal(msg)
    if err != nil {
        return fmt.Errorf("failed to marshal job execution message: %w", err)
    }

    // Create Kafka message with scheduler ID as key for partitioning
    kafkaMsg := kafka.Message{
        Key:   []byte(msg.SchedulerID),
        Value: data,
        Time:  time.Now(),
        Headers: []kafka.Header{
            {Key: "message_type", Value: []byte("job_execution")},
            {Key: "scheduler_id", Value: []byte(msg.SchedulerID)},
            {Key: "run_id", Value: []byte(msg.RunID)},
        },
    }

    // Publish message
    if err := p.writer.WriteMessages(ctx, kafkaMsg); err != nil {
        return fmt.Errorf("failed to publish job execution message: %w", err)
    }

    log.Printf("📨 Published job execution: runId=%s, schedulerId=%s", msg.RunID, msg.SchedulerID)
    return nil
}

func (p *Producer) Close() error {
    return p.writer.Close()
}

// Consumer for job execution messages
type Consumer struct {
    reader *kafka.Reader
}

func NewConsumer(brokers string, topic string, groupID string) *Consumer {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:     []string{brokers},
        Topic:       topic,
        GroupID:     groupID,
        StartOffset: kafka.FirstOffset, // Start from beginning for testing
        MinBytes:    10e3,              // 10KB
        MaxBytes:    10e6,              // 10MB
        MaxWait:     1 * time.Second,
    })

    log.Printf("✅ Kafka consumer initialized for topic: %s, group: %s", topic, groupID)
    return &Consumer{reader: reader}
}

func (c *Consumer) ReadMessage(ctx context.Context) (*JobExecutionMessage, error) {
    kafkaMsg, err := c.reader.ReadMessage(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to read Kafka message: %w", err)
    }

    var msg JobExecutionMessage
    if err := json.Unmarshal(kafkaMsg.Value, &msg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal job execution message: %w", err)
    }

    return &msg, nil
}

func (c *Consumer) Close() error {
    return c.reader.Close()
}
