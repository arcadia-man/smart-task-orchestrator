package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"smart-task-orchestrator/internal/config"
	"smart-task-orchestrator/internal/db"
	"smart-task-orchestrator/internal/jobs"
	"smart-task-orchestrator/internal/kafka"
	"smart-task-orchestrator/internal/retry"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB
	mongoDB, err := db.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoDB.Close()

	// Initialize services
	jobService := jobs.NewService(mongoDB.Database)
	consumer := kafka.NewConsumer(cfg.KafkaBroker, "jobs.execute", "worker-group-1")
	producer := kafka.NewProducer(cfg.KafkaBroker)
	defer consumer.Close()
	defer producer.Close()

	log.Println("🔄 Worker started, waiting for jobs...")

	ctx := context.Background()
	errorCount := 0
	maxErrors := 10

	for {
		// Read message from Kafka
		msg, err := consumer.ReadMessage(ctx)
		if err != nil {
			errorCount++
			log.Printf("Error reading message (%d/%d): %v", errorCount, maxErrors, err)

			// If too many consecutive errors, exit to prevent infinite loop
			if errorCount >= maxErrors {
				log.Fatal("Too many consecutive errors, shutting down worker")
			}

			// Exponential backoff for errors
			backoffTime := time.Duration(errorCount) * time.Second
			log.Printf("Waiting %v before retry...", backoffTime)
			time.Sleep(backoffTime)
			continue
		}

		// Reset error count on successful read
		errorCount = 0

		// Process the job
		processJob(ctx, jobService, producer, msg)
	}
}

func processJob(ctx context.Context, jobService *jobs.Service, producer *kafka.Producer, msg *kafka.JobMessage) {
	jobID := msg.JobID
	log.Printf("🔨 Processing job: %s", jobID)

	// Update job status to running
	err := jobService.UpdateJobStatus(ctx, jobID, jobs.StatusRunning, "Job execution started")
	if err != nil {
		log.Printf("Failed to update job status: %v", err)
		return
	}

	// Simulate job execution (replace with actual business logic)
	success := simulateJobExecution(msg.Payload)

	if success {
		// Job completed successfully
		err = jobService.UpdateJobStatus(ctx, jobID, jobs.StatusCompleted, "Job completed successfully")
		if err != nil {
			log.Printf("Failed to update job status: %v", err)
		}
		log.Printf("✅ Job %s completed successfully", jobID)
	} else {
		// Job failed, handle retry logic
		handleJobFailure(ctx, jobService, producer, jobID)
	}
}

func simulateJobExecution(payload map[string]any) bool {
	// Simulate processing time
	processingTime := time.Duration(rand.Intn(3)+1) * time.Second
	time.Sleep(processingTime)

	// Simulate 70% success rate
	return rand.Float32() < 0.7
}

func handleJobFailure(ctx context.Context, jobService *jobs.Service, producer *kafka.Producer, jobID string) {
	log.Printf("❌ Job %s failed", jobID)

	// Get current job to check retry count
	job, err := jobService.GetJobByID(ctx, jobID)
	if err != nil {
		log.Printf("Failed to get job: %v", err)
		return
	}

	// Check if we should retry
	if retry.ShouldRetry(job.RetryCount, job.MaxRetries) {
		// Calculate backoff delay
		delay := retry.CalculateBackoff(job.RetryCount)
		log.Printf("🔄 Scheduling retry for job %s after %v", jobID, delay)

		// Update job status
		err = jobService.UpdateJobStatus(ctx, jobID, jobs.StatusScheduled,
			fmt.Sprintf("Job failed, scheduled for retry after %v", delay))
		if err != nil {
			log.Printf("Failed to update job status: %v", err)
			return
		}

		// Schedule retry (in a real system, you might use a delay queue or scheduler)
		go func() {
			time.Sleep(delay)
			err := producer.PublishJob(ctx, "jobs.execute", jobID, job.Payload)
			if err != nil {
				log.Printf("Failed to republish job for retry: %v", err)
			} else {
				jobService.UpdateJobStatus(ctx, jobID, jobs.StatusQueued, "Job queued for retry")
			}
		}()
	} else {
		// Max retries exceeded, move to DLQ
		log.Printf("💀 Job %s exceeded max retries, moving to DLQ", jobID)

		err = jobService.UpdateJobStatus(ctx, jobID, jobs.StatusFailed, "Job failed after maximum retries")
		if err != nil {
			log.Printf("Failed to update job status: %v", err)
		}

		// Publish to DLQ
		err = producer.PublishJob(ctx, "jobs.failed", jobID, job.Payload)
		if err != nil {
			log.Printf("Failed to publish to DLQ: %v", err)
		}
	}
}
