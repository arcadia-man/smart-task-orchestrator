package models

import (
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type JobType string

const (
    JobTypeImmediate JobType = "immediate"
    JobTypeCron      JobType = "cron"
    JobTypeInterval  JobType = "interval"
)

type SchedulerStatus string

const (
    SchedulerStatusActive   SchedulerStatus = "active"
    SchedulerStatusInactive SchedulerStatus = "inactive"
    SchedulerStatusPaused   SchedulerStatus = "paused"
)

type SchedulerDefinition struct {
    ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name            string             `json:"name" bson:"name"`
    Description     string             `json:"description" bson:"description"`
    Image           string             `json:"image" bson:"image"`
    JobType         JobType            `json:"jobType" bson:"job_type"`
    CronExpr        string             `json:"cronExpr" bson:"cron_expr"`
    IntervalSeconds int                `json:"intervalSeconds" bson:"interval_seconds"`
    Command         string             `json:"command" bson:"command"`
    Status          SchedulerStatus    `json:"status" bson:"status"`
    Timezone        string             `json:"timezone" bson:"timezone"`
    Generation      int                `json:"generation" bson:"generation"`
    CreatedAt       time.Time          `json:"createdAt" bson:"created_at"`
    UpdatedAt       time.Time          `json:"updatedAt" bson:"updated_at"`
    CreatedBy       primitive.ObjectID `json:"createdBy" bson:"created_by"`
    UpdatedBy       primitive.ObjectID `json:"updatedBy" bson:"updated_by"`
}

type PrecomputeStatus string

const (
    PrecomputeStatusPending    PrecomputeStatus = "pending"
    PrecomputeStatusDispatched PrecomputeStatus = "dispatched"
    PrecomputeStatusCanceled   PrecomputeStatus = "canceled"
    PrecomputeStatusDiscarded  PrecomputeStatus = "discarded"
)

type SchedulerPrecompute struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    SchedulerID primitive.ObjectID `json:"schedulerId" bson:"scheduler_id"`
    RunAt       time.Time          `json:"runAt" bson:"run_at"`
    Generation  int                `json:"generation" bson:"generation"`
    Status      PrecomputeStatus   `json:"status" bson:"status"`
    CreatedAt   time.Time          `json:"createdAt" bson:"created_at"`
}

type ExecutionStatus string

const (
    ExecutionStatusPending   ExecutionStatus = "pending"
    ExecutionStatusRunning   ExecutionStatus = "running"
    ExecutionStatusSuccess   ExecutionStatus = "success"
    ExecutionStatusFailed    ExecutionStatus = "failed"
    ExecutionStatusDiscarded ExecutionStatus = "discarded"
)

type SchedulerHistory struct {
    ID           primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
    SchedulerID  primitive.ObjectID  `json:"schedulerId" bson:"scheduler_id"`
    PrecomputeID *primitive.ObjectID `json:"precomputeId" bson:"precompute_id"`
    RunID        primitive.ObjectID  `json:"runId" bson:"run_id"`
    ExecutedBy   primitive.ObjectID  `json:"executedBy" bson:"executed_by"`
    Status       ExecutionStatus     `json:"status" bson:"status"`
    StartTime    time.Time           `json:"startTime" bson:"start_time"`
    EndTime      *time.Time          `json:"endTime" bson:"end_time"`
    Command      string              `json:"command" bson:"command"`
    ExitCode     *int                `json:"exitCode" bson:"exit_code"`
    ProcessID    string              `json:"processId" bson:"process_id"`
    LogBlobID    string              `json:"logBlobId" bson:"log_blob_id"`
    LogText      string              `json:"logText" bson:"log_text"`
    ErrorMessage string              `json:"errorMessage" bson:"error_message"`
    CreatedAt    time.Time           `json:"createdAt" bson:"created_at"`
    UpdatedAt    time.Time           `json:"updatedAt" bson:"updated_at"`
}

type AuditLog struct {
    ID          primitive.ObjectID     `json:"id" bson:"_id,omitempty"`
    Action      string                 `json:"action" bson:"action"`
    EntityType  string                 `json:"entityType" bson:"entity_type"`
    EntityID    primitive.ObjectID     `json:"entityId" bson:"entity_id"`
    PerformedBy primitive.ObjectID     `json:"performedBy" bson:"performed_by"`
    PerformedAt time.Time              `json:"performedAt" bson:"performed_at"`
    Details     map[string]interface{} `json:"details" bson:"details"`
}
