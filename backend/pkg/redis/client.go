package redis

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

type Client struct {
    *redis.Client
    ShardCount int
}

func NewClient(url string, shardCount int) (*Client, error) {
    opts, err := redis.ParseURL(url)
    if err != nil {
        return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
    }

    rdb := redis.NewClient(opts)

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := rdb.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    client := &Client{
        Client:     rdb,
        ShardCount: shardCount,
    }

    log.Printf("✅ Connected to Redis with %d shards", shardCount)
    return client, nil
}

// GetShardKey returns the Redis key for a specific shard
func (c *Client) GetShardKey(shard int) string {
    return fmt.Sprintf("sched:zset:%d", shard)
}

// GetSchedulerShard returns the shard number for a scheduler ID
func (c *Client) GetSchedulerShard(schedulerID string) int {
    // Simple hash-based sharding
    hash := 0
    for _, char := range schedulerID {
        hash = hash*31 + int(char)
    }
    if hash < 0 {
        hash = -hash
    }
    return hash % c.ShardCount
}

// AddScheduledItem adds an item to the appropriate shard
func (c *Client) AddScheduledItem(ctx context.Context, schedulerID string, precomputeID string, runAtMs int64) error {
    shard := c.GetSchedulerShard(schedulerID)
    shardKey := c.GetShardKey(shard)

    member := fmt.Sprintf("precompute:%s", precomputeID)

    return c.ZAdd(ctx, shardKey, redis.Z{
        Score:  float64(runAtMs),
        Member: member,
    }).Err()
}

// PopDueItems pops all items from a shard that are due (score <= now)
func (c *Client) PopDueItems(ctx context.Context, shard int, nowMs int64) ([]string, error) {
    shardKey := c.GetShardKey(shard)

    // Use Lua script for atomic pop operation
    luaScript := `
        local key = KEYS[1]
        local maxScore = ARGV[1]
        
        -- Get items with score <= maxScore
        local items = redis.call('ZRANGEBYSCORE', key, 0, maxScore)
        
        if #items > 0 then
            -- Remove the items
            redis.call('ZREMRANGEBYSCORE', key, 0, maxScore)
        end
        
        return items
    `

    result, err := c.Eval(ctx, luaScript, []string{shardKey}, nowMs).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to pop due items: %w", err)
    }

    items, ok := result.([]interface{})
    if !ok {
        return nil, fmt.Errorf("unexpected result type from Redis")
    }

    stringItems := make([]string, len(items))
    for i, item := range items {
        stringItems[i] = item.(string)
    }

    return stringItems, nil
}

// RemoveSchedulerItems removes all items for a specific scheduler (for invalidation)
func (c *Client) RemoveSchedulerItems(ctx context.Context, schedulerID string, generation int) error {
    shard := c.GetSchedulerShard(schedulerID)
    shardKey := c.GetShardKey(shard)

    // Get all members and filter by scheduler (this could be optimized)
    members, err := c.ZRange(ctx, shardKey, 0, -1).Result()
    if err != nil {
        return fmt.Errorf("failed to get shard members: %w", err)
    }

    // Remove matching members (simplified - in production, include generation check)
    // This is a placeholder implementation - in production you'd parse member strings
    // to check scheduler ID and generation before removal
    _ = members // TODO: implement proper filtering and removal

    return nil
}

// GetShardStats returns statistics for a shard
func (c *Client) GetShardStats(ctx context.Context, shard int) (int64, error) {
    shardKey := c.GetShardKey(shard)
    return c.ZCard(ctx, shardKey).Result()
}
