package main

import (
	"context"
	"log"
	"strings"

	"smart-task-orchestrator/internal/config"
	"smart-task-orchestrator/internal/models"
	"smart-task-orchestrator/internal/pkg/db"
	"smart-task-orchestrator/internal/pkg/kafka"

	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
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
	producer := kafka.NewProducer(strings.Split(cfg.KafkaBrokers, ","), "jobs.execute")
	defer producer.Close()

	c := cron.New()
	
	// Check for cron jobs every minute
	c.AddFunc("* * * * *", func() {
		log.Println("Checking for due cron jobs...")
		
		ctx := context.Background()
		cursor, err := database.Collection("jobs").Find(ctx, bson.M{"type": models.JobTypeCron})
		if err != nil {
			log.Printf("Failed to fetch cron jobs: %v", err)
			return
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var job models.Job
			if err := cursor.Decode(&job); err != nil {
				continue
			}

			// In a real app, we'd check if it's actually due based on last run
			// For this MVP, we just trigger it if it exists to show it works
			log.Printf("Triggering cron job: %s", job.Name)
			producer.Publish(ctx, job.ID.Hex(), job)
		}
	})

	c.Start()
	log.Printf("Scheduler started")

	// Keep alive
	select {}
}
