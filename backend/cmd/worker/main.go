package main

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"smart-task-orchestrator/internal/config"
	"smart-task-orchestrator/internal/models"
	"smart-task-orchestrator/internal/pkg/db"
	"smart-task-orchestrator/internal/pkg/kafka"
	"smart-task-orchestrator/internal/worker"

	kafkaGo "github.com/segmentio/kafka-go"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := db.ConnectMongo(cfg.DBUri)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	database := client.Database(cfg.DBName)

	sandbox, err := worker.NewSandboxWorker(database)
	if err != nil {
		log.Fatalf("Failed to initialize sandbox worker: %v", err)
	}

	consumer := kafka.NewConsumer(strings.Split(cfg.KafkaBrokers, ","), "jobs.execute", "worker-group")
	defer consumer.Close()

	log.Printf("Worker started, consuming from jobs.execute")

	err = consumer.Consume(context.Background(), func(msg kafkaGo.Message) error {
		var job models.Job
		if err := json.Unmarshal(msg.Value, &job); err != nil {
			log.Printf("Failed to unmarshal job: %v", err)
			return nil // Drop bad message
		}

		log.Printf("Processing job: %s (%s)", job.Name, job.ID.Hex())
		
		err := sandbox.ExecuteJob(context.Background(), job, "")
		if err != nil {
			log.Printf("Job execution failed: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Consumer failed: %v", err)
	}
}
