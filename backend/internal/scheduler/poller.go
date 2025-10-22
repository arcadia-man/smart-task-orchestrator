package scheduler

import (
    "context"
    "fmt"
    "log"
    "strings"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"

    "smart-task-orchestrator/internal/models"
    "smart-task-orchestrator/pkg/kafka"
    "smart-task-orchestrator/pkg/redis"
)

type PollerService struct {
    db          *mongo.Database
    redisClient *redis.Client
    producer    *kafka.Producer
    shardID     int
    ticker      *time.Ticker
    stopChan    chan struct{}
}

func NewPollerService(db *mongo.Database, redisClient *redis.Client, producer *kafka.Producer, shardID int) *PollerService {
    return &PollerService{
        db:          db,
        redisClient: redisClient,
        producer:    producer,
        shardID:     shardID,
        stopChan:    make(chan struct{}),
    }
}

func (p *PollerService) Start() {
    log.Printf("🔄 Starting poller service for shard %d (every 1 second)", p.shardID)

    // Run every 1 second for precise timing
    p.ticker = time.NewTicker(1 * time.Second)

    go func() {
        for {
            select {
            case <-p.ticker.C:
                p.pollAndDispatch()
            case <-p.stopChan:
                return
            }
        }
    }()
}

func (p *PollerService) Stop() {
    log.Printf("🛑 Stopping poller service for shard %d", p.shardID)
    if p.ticker != nil {
        p.ticker.Stop()
    }
    close(p.stopChan)
}

func (p *PollerService) pollAndDispatch() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    now := time.Now()
    nowMs := now.UnixMilli()

    // Pop due items from Redis shard
    dueItems, err := p.redisClient.PopDueItems(ctx, p.shardID, nowMs)
    if err != nil {
        log.Printf("❌ Failed to pop due items from shard %d: %v", p.shardID, err)
        return
    }

    if len(dueItems) == 0 {
        return // No items due
    }

    log.Printf("⏰ Found %d due items in shard %d", len(dueItems), p.shardID)

    for _, item := range dueItems {
        if err := p.processItem(ctx, item, now); err != nil {
            log.Printf("❌ Failed to process item %s: %v", item, err)
        }
    }
}

func (p *PollerService) processItem(ctx context.Context, item string, now time.Time) error {
    // Parse item format: "precompute:ObjectID"
    parts := strings.Split(item, ":")
    if len(parts) != 2 || parts[0] != "precompute" {
        return fmt.Errorf("invalid item format: %s", item)
    }

    precomputeID, err := primitive.ObjectIDFromHex(parts[1])
    if err != nil {
        return fmt.Errorf("invalid precompute ID: %s", parts[1])
    }

    // Get precompute record from database
    precompute, err := p.getPrecomputeRecord(ctx, precomputeID)
    if err != nil {
        return fmt.Errorf("failed to get precompute record: %w", err)
    }

    // Get scheduler definition
    scheduler, err := p.getSchedulerDefinition(ctx, precompute.SchedulerID)
    if err != nil {
        return fmt.Errorf("failed to get scheduler definition: %w", err)
    }

    // Validate generation (prevent stale executions)
    if precompute.Generation != scheduler.Generation {
        log.Printf("⚠️  Discarding stale execution: precompute gen=%d, scheduler gen=%d",
            precompute.Generation, scheduler.Generation)

        return p.markPrecomputeDiscarded(ctx, precomputeID, "generation_mismatch")
    }

    // Check if scheduler is still active
    if scheduler.Status != models.SchedulerStatusActive {
        log.Printf("⚠️  Discarding execution for inactive scheduler: %s", scheduler.ID.Hex())
        return p.markPrecomputeDiscarded(ctx, precomputeID, "scheduler_inactive")
    }

    // Check for duplicate execution (safety check)
    if exists, err := p.checkDuplicateExecution(ctx, scheduler.ID, precompute.RunAt); err != nil {
        log.Printf("⚠️  Failed to check duplicate execution: %v", err)
    } else if exists {
        log.Printf("⚠️  Duplicate execution detected, skipping: %s at %v",
            scheduler.ID.Hex(), precompute.RunAt)
        return p.markPrecomputeDiscarded(ctx, precomputeID, "duplicate_execution")
    }

    // Create execution record
    runID := primitive.NewObjectID()
    historyRecord := models.SchedulerHistory{
        ID:           runID,
        SchedulerID:  scheduler.ID,
        PrecomputeID: &precomputeID,
        RunID:        runID,
        ExecutedBy:   primitive.NilObjectID, // System execution
        Status:       models.ExecutionStatusPending,
        StartTime:    now,
        Command:      scheduler.Command,
        CreatedAt:    now,
        UpdatedAt:    now,
    }

    if err := p.createHistoryRecord(ctx, historyRecord); err != nil {
        return fmt.Errorf("failed to create history record: %w", err)
    }

    // Mark precompute as dispatched
    if err := p.markPrecomputeDispatched(ctx, precomputeID); err != nil {
        log.Printf("⚠️  Failed to mark precompute as dispatched: %v", err)
        // Continue anyway - the execution should proceed
    }

    // Publish job execution message to Kafka
    jobMsg := kafka.JobExecutionMessage{
        RunID:       runID.Hex(),
        SchedulerID: scheduler.ID.Hex(),
        Generation:  scheduler.Generation,
        Image:       scheduler.Image,
        Command:     scheduler.Command,
        Env:         make(map[string]string), // TODO: Add environment variables
        TriggeredBy: "system",
        RunAt:       precompute.RunAt,
        Metadata: map[string]interface{}{
            "precompute_id":  precomputeID.Hex(),
            "scheduler_name": scheduler.Name,
            "job_type":       string(scheduler.JobType),
        },
    }

    if err := p.producer.PublishJobExecution(ctx, jobMsg); err != nil {
        // Mark history as failed if we can't publish
        p.updateHistoryStatus(ctx, runID, models.ExecutionStatusFailed,
            fmt.Sprintf("Failed to publish to Kafka: %v", err))
        return fmt.Errorf("failed to publish job execution: %w", err)
    }

    log.Printf("🚀 Dispatched job execution: runId=%s, scheduler=%s",
        runID.Hex(), scheduler.Name)

    return nil
}

func (p *PollerService) getPrecomputeRecord(ctx context.Context, precomputeID primitive.ObjectID) (*models.SchedulerPrecompute, error) {
    collection := p.db.Collection("scheduler_precompute")

    var precompute models.SchedulerPrecompute
    err := collection.FindOne(ctx, bson.M{"_id": precomputeID}).Decode(&precompute)
    if err != nil {
        return nil, fmt.Errorf("precompute record not found: %w", err)
    }

    return &precompute, nil
}

func (p *PollerService) getSchedulerDefinition(ctx context.Context, schedulerID primitive.ObjectID) (*models.SchedulerDefinition, error) {
    collection := p.db.Collection("scheduler_definition")

    var scheduler models.SchedulerDefinition
    err := collection.FindOne(ctx, bson.M{"_id": schedulerID}).Decode(&scheduler)
    if err != nil {
        return nil, fmt.Errorf("scheduler definition not found: %w", err)
    }

    return &scheduler, nil
}

func (p *PollerService) checkDuplicateExecution(ctx context.Context, schedulerID primitive.ObjectID, runAt time.Time) (bool, error) {
    collection := p.db.Collection("scheduler_history")

    // Check for existing execution within a small time window (±30 seconds)
    timeWindow := 30 * time.Second
    filter := bson.M{
        "scheduler_id": schedulerID,
        "start_time": bson.M{
            "$gte": runAt.Add(-timeWindow),
            "$lte": runAt.Add(timeWindow),
        },
        "status": bson.M{"$in": []string{
            string(models.ExecutionStatusPending),
            string(models.ExecutionStatusRunning),
            string(models.ExecutionStatusSuccess),
        }},
    }

    count, err := collection.CountDocuments(ctx, filter)
    if err != nil {
        return false, err
    }

    return count > 0, nil
}

func (p *PollerService) createHistoryRecord(ctx context.Context, record models.SchedulerHistory) error {
    collection := p.db.Collection("scheduler_history")
    _, err := collection.InsertOne(ctx, record)
    return err
}

func (p *PollerService) markPrecomputeDispatched(ctx context.Context, precomputeID primitive.ObjectID) error {
    collection := p.db.Collection("scheduler_precompute")

    update := bson.M{
        "$set": bson.M{
            "status": models.PrecomputeStatusDispatched,
        },
    }

    _, err := collection.UpdateOne(ctx, bson.M{"_id": precomputeID}, update)
    return err
}

func (p *PollerService) markPrecomputeDiscarded(ctx context.Context, precomputeID primitive.ObjectID, reason string) error {
    collection := p.db.Collection("scheduler_precompute")

    update := bson.M{
        "$set": bson.M{
            "status": models.PrecomputeStatusDiscarded,
        },
    }

    _, err := collection.UpdateOne(ctx, bson.M{"_id": precomputeID}, update)
    if err != nil {
        return err
    }

    log.Printf("🗑️  Discarded precompute %s: %s", precomputeID.Hex(), reason)
    return nil
}

func (p *PollerService) updateHistoryStatus(ctx context.Context, runID primitive.ObjectID, status models.ExecutionStatus, errorMsg string) error {
    collection := p.db.Collection("scheduler_history")

    update := bson.M{
        "$set": bson.M{
            "status":     status,
            "updated_at": time.Now(),
        },
    }

    if errorMsg != "" {
        update["$set"].(bson.M)["error_message"] = errorMsg
    }

    if status == models.ExecutionStatusSuccess || status == models.ExecutionStatusFailed {
        update["$set"].(bson.M)["end_time"] = time.Now()
    }

    _, err := collection.UpdateOne(ctx, bson.M{"_id": runID}, update)
    return err
}
