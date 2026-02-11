package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "smart-task-orchestrator/internal/config"
    "smart-task-orchestrator/internal/db"
    "smart-task-orchestrator/internal/worker"
    "smart-task-orchestrator/pkg/kafka"
    "smart-task-orchestrator/pkg/redis"
)

func main() {
    log.Println("🚀 Starting Smart Task Orchestrator - Worker Service")

    // Load configuration
    cfg := config.Load()

    // Initialize database
    mongoDB, err := db.NewMongoDB(cfg.MongoURI, cfg.DBName)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer mongoDB.Close()

    // Initialize Redis
    redisClient, err := redis.NewClient(cfg.RedisURL, 4)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer redisClient.Close()

    // Initialize Kafka consumer
    consumer := kafka.NewConsumer(cfg.KafkaBroker, "job_executions", "worker-group")
    defer consumer.Close()

    // Initialize worker service
    workerService := worker.NewWorkerService(mongoDB.Database, redisClient, consumer, cfg.DockerHost)

    // Start worker service
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        if err := workerService.Start(ctx); err != nil {
            log.Fatalf("Worker service failed: %v", err)
        }
    }()

    // Wait for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    log.Println("✅ Worker service started")
    log.Println("🔄 Listening for job execution messages...")
    log.Println("🐳 Docker host:", cfg.DockerHost)
    log.Println("📊 Press Ctrl+C to stop")

    <-sigChan

    log.Println("🛑 Shutting down worker service...")

    // Cancel context to stop worker
    cancel()

    // Give worker time to finish current jobs
    time.Sleep(5 * time.Second)

    log.Println("✅ Worker service stopped")
}
