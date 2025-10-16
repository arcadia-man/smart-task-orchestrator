package main

import (
	"context"
	"log"
	"time"

	"smart-task-orchestrator/internal/config"
	"smart-task-orchestrator/internal/db"
	"smart-task-orchestrator/internal/jobs"
	"smart-task-orchestrator/internal/kafka"

	"github.com/robfig/cron/v3"
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
	producer := kafka.NewProducer(cfg.KafkaBroker)
	defer producer.Close()

	// Initialize cron scheduler
	c := cron.New()

	// Schedule job to run every minute to check for scheduled jobs
	_, err = c.AddFunc("@every 1m", func() {
		processScheduledJobs(jobService, producer)
	})
	if err != nil {
		log.Fatal("Failed to add cron job:", err)
	}

	log.Println("⏰ Scheduler started, checking for scheduled jobs every minute...")
	c.Start()

	// Keep the scheduler running
	select {}
}

func processScheduledJobs(jobService *jobs.Service, producer *kafka.Producer) {
	ctx := context.Background()

	log.Println("🔍 Checking for scheduled jobs...")

	scheduledJobs, err := jobService.GetScheduledJobs(ctx)
	if err != nil {
		log.Printf("Failed to get scheduled jobs: %v", err)
		return
	}

	if len(scheduledJobs) == 0 {
		log.Println("No scheduled jobs found")
		return
	}

	log.Printf("Found %d scheduled jobs", len(scheduledJobs))

	for _, job := range scheduledJobs {
		log.Printf("📤 Publishing scheduled job: %s", job.ID.Hex())

		// Publish job to Kafka
		err := producer.PublishJob(ctx, "jobs.execute", job.ID.Hex(), job.Payload)
		if err != nil {
			log.Printf("Failed to publish job %s: %v", job.ID.Hex(), err)
			continue
		}

		// Update job status to queued
		err = jobService.UpdateJobStatus(ctx, job.ID.Hex(), jobs.StatusQueued, "Job queued by scheduler")
		if err != nil {
			log.Printf("Failed to update job status: %v", err)
		}

		// For cron jobs, schedule next run
		if job.Type == jobs.TypeCron && job.CronExpr != "" {
			scheduleNextRun(ctx, jobService, &job)
		}
	}
}

func scheduleNextRun(ctx context.Context, jobService *jobs.Service, job *jobs.Job) {
	// Parse cron expression and calculate next run time
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, err := parser.Parse(job.CronExpr)
	if err != nil {
		log.Printf("Invalid cron expression for job %s: %v", job.ID.Hex(), err)
		return
	}

	nextRun := schedule.Next(time.Now())
	log.Printf("Next run for job %s scheduled at: %v", job.ID.Hex(), nextRun)

	// For now, we'll just log the next scheduled time
	// In a production system, you would create a new job instance
	log.Printf("Next cron job instance for '%s' would be scheduled at: %v", job.Name, nextRun)
}
