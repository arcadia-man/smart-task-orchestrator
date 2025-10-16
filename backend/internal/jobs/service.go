package jobs

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	collection *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		collection: db.Collection("jobs"),
	}
}

func (s *Service) CreateJob(ctx context.Context, req CreateJobRequest) (*Job, error) {
	now := time.Now()

	job := &Job{
		Name:       req.Name,
		Type:       req.Type,
		Payload:    req.Payload,
		Status:     StatusScheduled,
		RetryCount: 0,
		MaxRetries: req.MaxRetries,
		CronExpr:   req.CronExpr,
		History: []HistoryEvent{
			{
				Timestamp: now,
				Event:     "created",
				Message:   "Job created successfully",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Set NextRunAt for immediate jobs
	if req.Type == TypeImmediate {
		job.NextRunAt = &now
	}

	result, err := s.collection.InsertOne(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	job.ID = result.InsertedID.(primitive.ObjectID)
	return job, nil
}

func (s *Service) GetAllJobs(ctx context.Context) ([]Job, error) {
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := s.collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, fmt.Errorf("failed to decode jobs: %w", err)
	}

	return jobs, nil
}

func (s *Service) GetJobByID(ctx context.Context, id string) (*Job, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid job ID: %w", err)
	}

	var job Job
	err = s.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&job)
	if err != nil {
		return nil, fmt.Errorf("job not found: %w", err)
	}

	return &job, nil
}

func (s *Service) UpdateJobStatus(ctx context.Context, id string, status JobStatus, message string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
		"$push": bson.M{
			"history": HistoryEvent{
				Timestamp: time.Now(),
				Event:     string(status),
				Message:   message,
			},
		},
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (s *Service) IncrementRetryCount(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid job ID: %w", err)
	}

	update := bson.M{
		"$inc": bson.M{"retryCount": 1},
		"$set": bson.M{"updatedAt": time.Now()},
		"$push": bson.M{
			"history": HistoryEvent{
				Timestamp: time.Now(),
				Event:     "retry",
				Message:   "Job retry attempt",
			},
		},
	}

	_, err = s.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (s *Service) GetScheduledJobs(ctx context.Context) ([]Job, error) {
	now := time.Now()
	filter := bson.M{
		"status": StatusScheduled,
		"nextRunAt": bson.M{
			"$lte": now,
		},
	}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scheduled jobs: %w", err)
	}
	defer cursor.Close(ctx)

	var jobs []Job
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, fmt.Errorf("failed to decode scheduled jobs: %w", err)
	}

	return jobs, nil
}
