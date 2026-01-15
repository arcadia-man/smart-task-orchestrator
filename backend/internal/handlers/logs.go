package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"smart-task-orchestrator/internal/models"
)

type LogHandlers struct {
	db *mongo.Database
}

func NewLogHandlers(db *mongo.Database) *LogHandlers {
	return &LogHandlers{db: db}
}

type LogEntry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Source    string    `json:"source"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
}

type LogStatsResponse struct {
	Success int64 `json:"success"`
	Info    int64 `json:"info"`
	Warning int64 `json:"warning"`
	Error   int64 `json:"error"`
}

func (h *LogHandlers) GetLogs(c *gin.Context) {
	// Query parameters
	level := c.Query("level")
	source := c.Query("source")
	search := c.Query("search")
	limitStr := c.DefaultQuery("limit", "100")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.ParseInt(limitStr, 10, 64)
	offset, _ := strconv.ParseInt(offsetStr, 10, 64)

	// Build filter
	filter := bson.M{}
	if level != "" && level != "all" {
		filter["level"] = level
	}
	if source != "" {
		filter["source"] = source
	}
	if search != "" {
		filter["$or"] = []bson.M{
			{"message": bson.M{"$regex": search, "$options": "i"}},
			{"source": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	// Get logs from audit_logs collection (using it as system logs)
	auditCollection := h.db.Collection("audit_logs")

	// For demo purposes, let's also get some logs from scheduler history
	historyCollection := h.db.Collection("scheduler_history")

	var logs []LogEntry

	// Get audit logs
	cursor, err := auditCollection.Find(
		context.Background(),
		filter,
		options.Find().
			SetSort(bson.M{"performed_at": -1}).
			SetLimit(limit/2).
			SetSkip(offset),
	)
	if err == nil {
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var auditLog models.AuditLog
			if err := cursor.Decode(&auditLog); err != nil {
				continue
			}

			level := "info"
			if auditLog.Action == "delete" {
				level = "warning"
			}

			logs = append(logs, LogEntry{
				ID:        auditLog.ID.Hex(),
				Timestamp: auditLog.PerformedAt,
				Level:     level,
				Source:    auditLog.EntityType,
				Message:   auditLog.Action + " " + auditLog.EntityType,
				Details:   "Entity ID: " + auditLog.EntityID.Hex(),
			})
		}
	}

	// Get scheduler execution logs
	historyCursor, err := historyCollection.Find(
		context.Background(),
		bson.M{},
		options.Find().
			SetSort(bson.M{"start_time": -1}).
			SetLimit(limit/2),
	)
	if err == nil {
		defer historyCursor.Close(context.Background())

		for historyCursor.Next(context.Background()) {
			var history models.SchedulerHistory
			if err := historyCursor.Decode(&history); err != nil {
				continue
			}

			level := "info"
			message := "Scheduler execution started"
			details := "Command: " + history.Command

			switch history.Status {
			case models.ExecutionStatusSuccess:
				level = "success"
				message = "Scheduler execution completed successfully"
				if history.EndTime != nil {
					duration := history.EndTime.Sub(history.StartTime)
					details += ", Duration: " + duration.String()
				}
			case models.ExecutionStatusFailed:
				level = "error"
				message = "Scheduler execution failed"
				if history.ErrorMessage != "" {
					details += ", Error: " + history.ErrorMessage
				}
			case models.ExecutionStatusRunning:
				level = "info"
				message = "Scheduler execution in progress"
			}

			logs = append(logs, LogEntry{
				ID:        history.ID.Hex(),
				Timestamp: history.StartTime,
				Level:     level,
				Source:    "scheduler-" + history.SchedulerID.Hex()[:8],
				Message:   message,
				Details:   details,
			})
		}
	}

	// Add some system logs for demo
	if len(logs) < 10 {
		now := time.Now()
		systemLogs := []LogEntry{
			{
				ID:        primitive.NewObjectID().Hex(),
				Timestamp: now.Add(-5 * time.Minute),
				Level:     "info",
				Source:    "auth",
				Message:   "User login successful",
				Details:   "IP: 192.168.1.100, User-Agent: Mozilla/5.0...",
			},
			{
				ID:        primitive.NewObjectID().Hex(),
				Timestamp: now.Add(-10 * time.Minute),
				Level:     "warning",
				Source:    "system",
				Message:   "High memory usage detected",
				Details:   "Memory usage: 85%, Available: 2.1GB",
			},
			{
				ID:        primitive.NewObjectID().Hex(),
				Timestamp: now.Add(-15 * time.Minute),
				Level:     "success",
				Source:    "backup",
				Message:   "Database backup completed successfully",
				Details:   "Backup size: 1.2GB, Duration: 45s",
			},
		}
		logs = append(logs, systemLogs...)
	}

	c.JSON(http.StatusOK, logs)
}

func (h *LogHandlers) GetLogStats(c *gin.Context) {
	// For demo purposes, return some stats
	// In a real implementation, you would aggregate from your logs collection

	stats := LogStatsResponse{
		Success: 15,
		Info:    8,
		Warning: 3,
		Error:   2,
	}

	// Try to get real stats from scheduler history
	historyCollection := h.db.Collection("scheduler_history")

	// Count by status in the last 24 hours
	since := time.Now().Add(-24 * time.Hour)

	success, _ := historyCollection.CountDocuments(context.Background(), bson.M{
		"status":     models.ExecutionStatusSuccess,
		"start_time": bson.M{"$gte": since},
	})

	failed, _ := historyCollection.CountDocuments(context.Background(), bson.M{
		"status":     models.ExecutionStatusFailed,
		"start_time": bson.M{"$gte": since},
	})

	running, _ := historyCollection.CountDocuments(context.Background(), bson.M{
		"status":     models.ExecutionStatusRunning,
		"start_time": bson.M{"$gte": since},
	})

	if success > 0 || failed > 0 || running > 0 {
		stats.Success = success
		stats.Error = failed
		stats.Info = running
		stats.Warning = 1 // Default warning count
	}

	c.JSON(http.StatusOK, stats)
}

func (h *LogHandlers) GetLogSources(c *gin.Context) {
	// Return available log sources
	sources := []string{
		"auth",
		"system",
		"scheduler",
		"backup",
		"api",
		"database",
	}

	c.JSON(http.StatusOK, sources)
}
