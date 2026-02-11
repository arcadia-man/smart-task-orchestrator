package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "smart-task-orchestrator/internal/config"
    "smart-task-orchestrator/internal/db"
    "smart-task-orchestrator/internal/scheduler"
    "smart-task-orchestrator/pkg/kafka"
    "smart-task-orchestrator/pkg/redis"
)

func main() {
    log.Println("🚀 Starting Smart Task Orchestrator - Scheduler Service")

    // Load configuration
    cfg := config.Load()

    // Initialize database
    mongoDB, err := db.NewMongoDB(cfg.MongoURI, cfg.DBName)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer mongoDB.Close()

    // Initialize Redis with 4 shards
    redisClient, err := redis.NewClient(cfg.RedisURL, 4)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer redisClient.Close()

    // Initialize Kafka producer
    producer := kafka.NewProducer(cfg.KafkaBroker, "job_executions")
    defer producer.Close()

    // Initialize services
    precomputeService := scheduler.NewPrecomputeService(mongoDB.Database, redisClient)

    // For now, run poller for shard 0 only
    // In production, you'd run multiple scheduler instances with different shard IDs
    pollerService := scheduler.NewPollerService(mongoDB.Database, redisClient, producer, 0)

    // Start services
    precomputeService.Start()
    pollerService.Start()

    // Wait for shutdown signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    log.Println("✅ Scheduler service started")
    log.Println("🔄 Precompute service: Running every 5 minutes")
    log.Println("⏰ Poller service: Running every 1 second for shard 0")
    log.Println("📊 Press Ctrl+C to stop")

    <-sigChan

    log.Println("🛑 Shutting down scheduler service...")

    // Stop services gracefully
    precomputeService.Stop()
    pollerService.Stop()

    // Give services time to finish current operations
    time.Sleep(2 * time.Second)

    log.Println("✅ Scheduler service stopped")
}
