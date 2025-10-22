package scheduler

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/robfig/cron/v3"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "smart-task-orchestrator/internal/models"
    "smart-task-orchestrator/pkg/redis"
)

type PrecomputeService struct {
    db          *mongo.Database
    redisClient *redis.Client
    ticker      *time.Ticker
    stopChan    chan struct{}
}

func NewPrecomputeService(db *mongo.Database, redisClient *redis.Client) *PrecomputeService {
    return &PrecomputeService{
        db:          db,
        redisClient: redisClient,
        stopChan:    make(chan struct{}),
    }
}

func (p *PrecomputeService) Start() {
    log.Println("🔄 Starting precompute service (every 5 minutes)")

    // Run immediately on start
    go p.runPrecompute()

    // Then run every 5 minutes
    p.ticker = time.NewTicker(5 * time.Minute)

    go func() {
        for {
            select {
            case <-p.ticker.C:
                p.runPrecompute()
            case <-p.stopChan:
                return
            }
        }
    }()
}

func (p *PrecomputeService) Stop() {
    log.Println("🛑 Stopping precompute service")
    if p.ticker != nil {
        p.ticker.Stop()
    }
    close(p.stopChan)
}

func (p *PrecomputeService) runPrecompute() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
    defer cancel()

    log.Println("⚙️  Running precompute cycle...")

    // Get all active schedulers
    schedulers, err := p.getActiveSchedulers(ctx)
    if err != nil {
        log.Printf("❌ Failed to get active schedulers: %v", err)
        return
    }

    log.Printf("📋 Processing %d active schedulers", len(schedulers))

    for _, scheduler := range schedulers {
        if err := p.precomputeScheduler(ctx, scheduler); err != nil {
            log.Printf("❌ Failed to precompute scheduler %s: %v", scheduler.ID.Hex(), err)
        }
    }

    log.Println("✅ Precompute cycle completed")
}

func (p *PrecomputeService) getActiveSchedulers(ctx context.Context) ([]models.SchedulerDefinition, error) {
    collection := p.db.Collection("scheduler_definition")

    filter := bson.M{
        "status":   models.SchedulerStatusActive,
        "job_type": bson.M{"$in": []string{string(models.JobTypeCron), string(models.JobTypeInterval)}},
    }

    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, fmt.Errorf("failed to find active schedulers: %w", err)
    }
    defer cursor.Close(ctx)

    var schedulers []models.SchedulerDefinition
    if err := cursor.All(ctx, &schedulers); err != nil {
        return nil, fmt.Errorf("failed to decode schedulers: %w", err)
    }

    return schedulers, nil
}

func (p *PrecomputeService) precomputeScheduler(ctx context.Context, scheduler models.SchedulerDefinition) error {
    now := time.Now()
    lookAhead := 15 * time.Minute // 15-minute lookahead window
    endTime := now.Add(lookAhead)

    var nextRuns []time.Time
    var err error

    switch scheduler.JobType {
    case models.JobTypeCron:
        nextRuns, err = p.calculateCronRuns(scheduler.CronExpr, scheduler.Timezone, now, endTime)
    case models.JobTypeInterval:
        nextRuns, err = p.calculateIntervalRuns(scheduler.IntervalSeconds, now, endTime)
    default:
        return fmt.Errorf("unsupported job type: %s", scheduler.JobType)
    }

    if err != nil {
        return fmt.Errorf("failed to calculate next runs: %w", err)
    }

    // Insert precomputed runs
    for _, runTime := range nextRuns {
        if err := p.insertPrecomputeRun(ctx, scheduler.ID, runTime, scheduler.Generation); err != nil {
            log.Printf("⚠️  Failed to insert precompute run for %s at %v: %v",
                scheduler.ID.Hex(), runTime, err)
        }
    }

    log.Printf("📅 Precomputed %d runs for scheduler %s", len(nextRuns), scheduler.Name)
    return nil
}

func (p *PrecomputeService) calculateCronRuns(cronExpr, timezone string, start, end time.Time) ([]time.Time, error) {
    // Parse timezone
    loc, err := time.LoadLocation(timezone)
    if err != nil {
        loc = time.UTC // Fallback to UTC
    }

    // Parse cron expression
    parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
    schedule, err := parser.Parse(cronExpr)
    if err != nil {
        return nil, fmt.Errorf("invalid cron expression: %w", err)
    }

    var runs []time.Time
    current := start.In(loc)

    for current.Before(end) {
        next := schedule.Next(current)
        if next.After(end) {
            break
        }
        runs = append(runs, next.UTC())
        current = next
    }

    return runs, nil
}

func (p *PrecomputeService) calculateIntervalRuns(intervalSeconds int, start, end time.Time) ([]time.Time, error) {
    if intervalSeconds <= 0 {
        return nil, fmt.Errorf("invalid interval: %d", intervalSeconds)
    }

    interval := time.Duration(intervalSeconds) * time.Second
    var runs []time.Time

    // Start from the next interval boundary
    current := start.Add(interval)

    for current.Before(end) {
        runs = append(runs, current)
        current = current.Add(interval)
    }

    return runs, nil
}

func (p *PrecomputeService) insertPrecomputeRun(ctx context.Context, schedulerID primitive.ObjectID, runAt time.Time, generation int) error {
    precomputeCollection := p.db.Collection("scheduler_precompute")

    precompute := models.SchedulerPrecompute{
        ID:          primitive.NewObjectID(),
        SchedulerID: schedulerID,
        RunAt:       runAt,
        Generation:  generation,
        Status:      models.PrecomputeStatusPending,
        CreatedAt:   time.Now(),
    }

    // Insert with upsert to handle duplicates
    filter := bson.M{
        "scheduler_id": schedulerID,
        "run_at":       runAt,
        "generation":   generation,
    }

    update := bson.M{
        "$setOnInsert": precompute,
    }

    opts := options.Update().SetUpsert(true)
    result, err := precomputeCollection.UpdateOne(ctx, filter, update, opts)
    if err != nil {
        return fmt.Errorf("failed to insert precompute run: %w", err)
    }

    // If this was an insert (not an update), add to Redis
    if result.UpsertedCount > 0 {
        runAtMs := runAt.UnixMilli()
        precomputeID := precompute.ID.Hex()

        if result.UpsertedID != nil {
            if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
                precomputeID = oid.Hex()
            }
        }

        if err := p.redisClient.AddScheduledItem(ctx, schedulerID.Hex(), precomputeID, runAtMs); err != nil {
            log.Printf("⚠️  Failed to add item to Redis: %v", err)
            // Don't return error - DB insert succeeded, Redis is just cache
        }
    }

    return nil
}
