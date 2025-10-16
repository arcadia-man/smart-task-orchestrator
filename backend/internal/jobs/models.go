package jobs

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobStatus string

const (
	StatusScheduled JobStatus = "scheduled"
	StatusQueued    JobStatus = "queued"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type JobType string

const (
	TypeImmediate JobType = "immediate"
	TypeCron      JobType = "cron"
)

type HistoryEvent struct {
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	Event     string    `json:"event" bson:"event"`
	Message   string    `json:"message" bson:"message"`
}

type Job struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Type       JobType            `json:"type" bson:"type"`
	Payload    map[string]any     `json:"payload" bson:"payload"`
	Status     JobStatus          `json:"status" bson:"status"`
	RetryCount int                `json:"retryCount" bson:"retryCount"`
	MaxRetries int                `json:"maxRetries" bson:"maxRetries"`
	NextRunAt  *time.Time         `json:"nextRunAt" bson:"nextRunAt,omitempty"`
	CronExpr   string             `json:"cronExpr,omitempty" bson:"cronExpr,omitempty"`
	History    []HistoryEvent     `json:"history" bson:"history"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type CreateJobRequest struct {
	Name       string         `json:"name" binding:"required"`
	Type       JobType        `json:"type" binding:"required"`
	Payload    map[string]any `json:"payload"`
	MaxRetries int            `json:"maxRetries"`
	CronExpr   string         `json:"cronExpr,omitempty"`
}
