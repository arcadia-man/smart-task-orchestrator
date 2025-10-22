package worker

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "smart-task-orchestrator/internal/models"
    "smart-task-orchestrator/pkg/docker"
    "smart-task-orchestrator/pkg/kafka"
    redisClient "smart-task-orchestrator/pkg/redis"
)

type WorkerService struct {
    db                *mongo.Database
    redisClient       *redisClient.Client
    consumer          *kafka.Consumer
    dockerClient      *docker.Client
    maxConcurrentJobs int
    jobSemaphore      chan struct{}
}

func NewWorkerService(db *mongo.Database, redisClient *redisClient.Client, consumer *kafka.Consumer, dockerHost string) *WorkerService {
    dockerClient, err := docker.NewClient(dockerHost)
    if err != nil {
        log.Fatalf("Failed to create Docker client: %v", err)
    }

    maxConcurrentJobs := 5 // Configurable

    return &WorkerService{
        db:                db,
        redisClient:       redisClient,
        consumer:          consumer,
        dockerClient:      dockerClient,
        maxConcurrentJobs: maxConcurrentJobs,
        jobSemaphore:      make(chan struct{}, maxConcurrentJobs),
    }
}

func (w *WorkerService) Start(ctx context.Context) error {
    log.Printf("🔄 Worker service started with max %d concurrent jobs", w.maxConcurrentJobs)

    for {
        select {
        case <-ctx.Done():
            log.Println("🛑 Worker service context cancelled")
            return nil
        default:
            // Read message from Kafka
            msg, err := w.consumer.ReadMessage(ctx)
            if err != nil {
                log.Printf("❌ Failed to read Kafka message: %v", err)
                time.Sleep(1 * time.Second)
                continue
            }

            // Process job in goroutine with semaphore
            go w.processJob(ctx, msg)
        }
    }
}

func (w *WorkerService) processJob(ctx context.Context, msg *kafka.JobExecutionMessage) {
    // Acquire semaphore slot
    select {
    case w.jobSemaphore <- struct{}{}:
        defer func() { <-w.jobSemaphore }()
    case <-ctx.Done():
        return
    }

    log.Printf("🔄 Processing job: runId=%s, scheduler=%s", msg.RunID, msg.SchedulerID)

    runID, err := primitive.ObjectIDFromHex(msg.RunID)
    if err != nil {
        log.Printf("❌ Invalid run ID: %s", msg.RunID)
        return
    }

    // Update history status to running
    if err := w.updateHistoryStatus(ctx, runID, models.ExecutionStatusRunning, ""); err != nil {
        log.Printf("❌ Failed to update history status: %v", err)
        return
    }

    // Execute the job
    result := w.executeJob(ctx, msg)

    // Update final status
    status := models.ExecutionStatusSuccess
    if result.ExitCode != 0 {
        status = models.ExecutionStatusFailed
    }

    if err := w.updateHistoryWithResult(ctx, runID, status, result); err != nil {
        log.Printf("❌ Failed to update history with result: %v", err)
    }

    // Stream final log message
    w.streamFinalLogMessage(ctx, msg.RunID, status, result.ExitCode)

    log.Printf("✅ Job completed: runId=%s, status=%s, exitCode=%d",
        msg.RunID, status, result.ExitCode)
}

func (w *WorkerService) executeJob(ctx context.Context, msg *kafka.JobExecutionMessage) *docker.ExecutionResult {
    // Create container name
    containerName := fmt.Sprintf("sched-%s-run-%s", msg.SchedulerID, msg.RunID)

    // Execute job with Docker
    result, err := w.dockerClient.ExecuteJob(ctx, &docker.JobSpec{
        Image:         msg.Image,
        Command:       msg.Command,
        ContainerName: containerName,
        Env:           msg.Env,
        RunID:         msg.RunID,
    })

    if err != nil {
        log.Printf("❌ Docker execution failed: %v", err)
        return &docker.ExecutionResult{
            ExitCode:     -1,
            ErrorMessage: err.Error(),
            Logs:         []string{fmt.Sprintf("Docker execution failed: %v", err)},
        }
    }

    // Stream logs to Redis for real-time viewing
    go w.streamLogsToRedis(ctx, msg.RunID, result.LogStream)

    return result
}

func (w *WorkerService) streamLogsToRedis(ctx context.Context, runID string, logStream <-chan string) {
    streamKey := fmt.Sprintf("logs:%s", runID)

    for {
        select {
        case logLine, ok := <-logStream:
            if !ok {
                return // Channel closed
            }

            // Add log line to Redis stream
            if err := w.redisClient.Client.XAdd(ctx, &redis.XAddArgs{
                Stream: streamKey,
                Values: map[string]interface{}{
                    "type":      "log",
                    "line":      logLine,
                    "timestamp": time.Now().UnixMilli(),
                },
            }).Err(); err != nil {
                log.Printf("⚠️  Failed to stream log to Redis: %v", err)
            }

        case <-ctx.Done():
            return
        }
    }
}

func (w *WorkerService) streamFinalLogMessage(ctx context.Context, runID string, status models.ExecutionStatus, exitCode int) {
    streamKey := fmt.Sprintf("logs:%s", runID)

    statusText := "success"
    if status == models.ExecutionStatusFailed {
        statusText = "failed"
    }

    if err := w.redisClient.Client.XAdd(ctx, &redis.XAddArgs{
        Stream: streamKey,
        Values: map[string]interface{}{
            "type":      "end",
            "status":    statusText,
            "exit_code": exitCode,
            "timestamp": time.Now().UnixMilli(),
        },
    }).Err(); err != nil {
        log.Printf("⚠️  Failed to stream final message to Redis: %v", err)
    }
}

func (w *WorkerService) updateHistoryStatus(ctx context.Context, runID primitive.ObjectID, status models.ExecutionStatus, errorMsg string) error {
    collection := w.db.Collection("scheduler_history")

    update := bson.M{
        "$set": bson.M{
            "status":     status,
            "updated_at": time.Now(),
        },
    }

    if errorMsg != "" {
        update["$set"].(bson.M)["error_message"] = errorMsg
    }

    if status == models.ExecutionStatusRunning {
        update["$set"].(bson.M)["start_time"] = time.Now()
    }

    _, err := collection.UpdateOne(ctx, bson.M{"_id": runID}, update)
    return err
}

func (w *WorkerService) updateHistoryWithResult(ctx context.Context, runID primitive.ObjectID, status models.ExecutionStatus, result *docker.ExecutionResult) error {
    collection := w.db.Collection("scheduler_history")

    // Combine all logs into a single text
    logText := ""
    for _, line := range result.Logs {
        logText += line + "\n"
    }

    update := bson.M{
        "$set": bson.M{
            "status":        status,
            "end_time":      time.Now(),
            "updated_at":    time.Now(),
            "exit_code":     result.ExitCode,
            "process_id":    result.ContainerID,
            "log_text":      logText,
            "error_message": result.ErrorMessage,
        },
    }

    _, err := collection.UpdateOne(ctx, bson.M{"_id": runID}, update)
    return err
}
