package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	Phone        string             `bson:"phone" json:"phone"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	APIKey       string             `bson:"api_key" json:"api_key"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

type JobType string

const (
	JobTypeOneTime JobType = "one-time"
	JobTypeCron    JobType = "cron"
	JobTypeSandbox JobType = "sandbox"
)

type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

type SandboxConfig struct {
	MinContainers int `bson:"min_containers" json:"min_containers"`
	MaxContainers int `bson:"max_containers" json:"max_containers"`
}

type Job struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         primitive.ObjectID `bson:"user_id" json:"user_id"`
	Name           string             `bson:"name" json:"name"`
	Type           JobType            `bson:"type" json:"type"`
	Image          string             `bson:"image" json:"image"`
	Command        string             `bson:"command" json:"command"`
	CronExpr       string             `bson:"cron_expr,omitempty" json:"cron_expr,omitempty"`
	Scaling        *SandboxConfig     `bson:"scaling,omitempty" json:"scaling,omitempty"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type Execution struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	JobID      primitive.ObjectID `bson:"job_id" json:"job_id"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status     JobStatus          `bson:"status" json:"status"`
	Logs       string             `bson:"logs" json:"logs"`
	InputData  string             `bson:"input_data,omitempty" json:"input_data,omitempty"`
	OutputData string             `bson:"output_data,omitempty" json:"output_data,omitempty"`
	StartedAt  time.Time          `bson:"started_at" json:"started_at"`
	EndedAt    *time.Time         `bson:"ended_at,omitempty" json:"ended_at,omitempty"`
}
