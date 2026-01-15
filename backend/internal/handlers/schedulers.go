package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"smart-task-orchestrator/internal/auth"
	"smart-task-orchestrator/internal/models"
)

type SchedulerHandlers struct {
	db *mongo.Database
}

func NewSchedulerHandlers(db *mongo.Database) *SchedulerHandlers {
	return &SchedulerHandlers{db: db}
}

type CreateSchedulerRequest struct {
	Name            string         `json:"name" binding:"required"`
	Description     string         `json:"description"`
	Image           string         `json:"image" binding:"required"`
	JobType         models.JobType `json:"jobType" binding:"required"`
	CronExpr        string         `json:"cronExpr"`
	IntervalSeconds int            `json:"intervalSeconds"`
	Command         string         `json:"command" binding:"required"`
	Timezone        string         `json:"timezone"`
}

type UpdateSchedulerRequest struct {
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Image           string                 `json:"image"`
	JobType         models.JobType         `json:"jobType"`
	CronExpr        string                 `json:"cronExpr"`
	IntervalSeconds int                    `json:"intervalSeconds"`
	Command         string                 `json:"command"`
	Status          models.SchedulerStatus `json:"status"`
	Timezone        string                 `json:"timezone"`
}

type SchedulerResponse struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	Image           string                 `json:"image"`
	JobType         models.JobType         `json:"jobType"`
	CronExpr        string                 `json:"cronExpr"`
	IntervalSeconds int                    `json:"intervalSeconds"`
	Command         string                 `json:"command"`
	Status          models.SchedulerStatus `json:"status"`
	Timezone        string                 `json:"timezone"`
	Generation      int                    `json:"generation"`
	LastRun         *time.Time             `json:"lastRun"`
	NextRun         *time.Time             `json:"nextRun"`
	LastStatus      string                 `json:"lastStatus"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
}

type SchedulerStatsResponse struct {
	Total    int64 `json:"total"`
	Active   int64 `json:"active"`
	Paused   int64 `json:"paused"`
	Inactive int64 `json:"inactive"`
}

func (h *SchedulerHandlers) GetSchedulers(c *gin.Context) {
	schedulersCollection := h.db.Collection("scheduler_definitions")
	historyCollection := h.db.Collection("scheduler_history")
	precomputeCollection := h.db.Collection("scheduler_precompute")

	// Get all schedulers
	cursor, err := schedulersCollection.Find(context.Background(), bson.M{}, options.Find().SetSort(bson.M{"created_at": -1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedulers"})
		return
	}
	defer cursor.Close(context.Background())

	var schedulers []SchedulerResponse
	for cursor.Next(context.Background()) {
		var scheduler models.SchedulerDefinition
		if err := cursor.Decode(&scheduler); err != nil {
			continue
		}

		// Get last execution
		var lastHistory models.SchedulerHistory
		historyCollection.FindOne(
			context.Background(),
			bson.M{"scheduler_id": scheduler.ID},
			options.FindOne().SetSort(bson.M{"start_time": -1}),
		).Decode(&lastHistory)

		// Get next run from precompute
		var nextPrecompute models.SchedulerPrecompute
		precomputeCollection.FindOne(
			context.Background(),
			bson.M{
				"scheduler_id": scheduler.ID,
				"status":       models.PrecomputeStatusPending,
			},
			options.FindOne().SetSort(bson.M{"run_at": 1}),
		).Decode(&nextPrecompute)

		var lastRun *time.Time
		var lastStatus string
		if !lastHistory.ID.IsZero() {
			lastRun = &lastHistory.StartTime
			lastStatus = string(lastHistory.Status)
		}

		var nextRun *time.Time
		if !nextPrecompute.ID.IsZero() {
			nextRun = &nextPrecompute.RunAt
		}

		schedulers = append(schedulers, SchedulerResponse{
			ID:              scheduler.ID.Hex(),
			Name:            scheduler.Name,
			Description:     scheduler.Description,
			Image:           scheduler.Image,
			JobType:         scheduler.JobType,
			CronExpr:        scheduler.CronExpr,
			IntervalSeconds: scheduler.IntervalSeconds,
			Command:         scheduler.Command,
			Status:          scheduler.Status,
			Timezone:        scheduler.Timezone,
			Generation:      scheduler.Generation,
			LastRun:         lastRun,
			NextRun:         nextRun,
			LastStatus:      lastStatus,
			CreatedAt:       scheduler.CreatedAt,
			UpdatedAt:       scheduler.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, schedulers)
}

func (h *SchedulerHandlers) GetSchedulerStats(c *gin.Context) {
	schedulersCollection := h.db.Collection("scheduler_definitions")

	total, _ := schedulersCollection.CountDocuments(context.Background(), bson.M{})
	active, _ := schedulersCollection.CountDocuments(context.Background(), bson.M{"status": models.SchedulerStatusActive})
	paused, _ := schedulersCollection.CountDocuments(context.Background(), bson.M{"status": models.SchedulerStatusPaused})
	inactive, _ := schedulersCollection.CountDocuments(context.Background(), bson.M{"status": models.SchedulerStatusInactive})

	c.JSON(http.StatusOK, SchedulerStatsResponse{
		Total:    total,
		Active:   active,
		Paused:   paused,
		Inactive: inactive,
	})
}

func (h *SchedulerHandlers) GetScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(schedulerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduler ID"})
		return
	}

	schedulersCollection := h.db.Collection("scheduler_definitions")
	var scheduler models.SchedulerDefinition
	err = schedulersCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&scheduler)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scheduler not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheduler"})
		return
	}

	c.JSON(http.StatusOK, SchedulerResponse{
		ID:              scheduler.ID.Hex(),
		Name:            scheduler.Name,
		Description:     scheduler.Description,
		Image:           scheduler.Image,
		JobType:         scheduler.JobType,
		CronExpr:        scheduler.CronExpr,
		IntervalSeconds: scheduler.IntervalSeconds,
		Command:         scheduler.Command,
		Status:          scheduler.Status,
		Timezone:        scheduler.Timezone,
		Generation:      scheduler.Generation,
		CreatedAt:       scheduler.CreatedAt,
		UpdatedAt:       scheduler.UpdatedAt,
	})
}

func (h *SchedulerHandlers) CreateScheduler(c *gin.Context) {
	var req CreateSchedulerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get current user from context
	userCtx, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate job type specific fields
	if req.JobType == models.JobTypeCron && req.CronExpr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cron expression is required for cron jobs"})
		return
	}
	if req.JobType == models.JobTypeInterval && req.IntervalSeconds <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Interval seconds must be greater than 0 for interval jobs"})
		return
	}

	// Set default timezone if not provided
	if req.Timezone == "" {
		req.Timezone = "UTC"
	}

	// Create scheduler
	now := time.Now()
	scheduler := models.SchedulerDefinition{
		Name:            req.Name,
		Description:     req.Description,
		Image:           req.Image,
		JobType:         req.JobType,
		CronExpr:        req.CronExpr,
		IntervalSeconds: req.IntervalSeconds,
		Command:         req.Command,
		Status:          models.SchedulerStatusActive,
		Timezone:        req.Timezone,
		Generation:      1,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       userCtx.UserID,
		UpdatedBy:       userCtx.UserID,
	}

	schedulersCollection := h.db.Collection("scheduler_definitions")
	result, err := schedulersCollection.InsertOne(context.Background(), scheduler)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create scheduler"})
		return
	}

	scheduler.ID = result.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusCreated, SchedulerResponse{
		ID:              scheduler.ID.Hex(),
		Name:            scheduler.Name,
		Description:     scheduler.Description,
		Image:           scheduler.Image,
		JobType:         scheduler.JobType,
		CronExpr:        scheduler.CronExpr,
		IntervalSeconds: scheduler.IntervalSeconds,
		Command:         scheduler.Command,
		Status:          scheduler.Status,
		Timezone:        scheduler.Timezone,
		Generation:      scheduler.Generation,
		CreatedAt:       scheduler.CreatedAt,
		UpdatedAt:       scheduler.UpdatedAt,
	})
}

func (h *SchedulerHandlers) UpdateScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(schedulerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduler ID"})
		return
	}

	var req UpdateSchedulerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Get current user from context
	userCtx, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	schedulersCollection := h.db.Collection("scheduler_definitions")

	// Build update document
	updateDoc := bson.M{
		"updated_at": time.Now(),
		"updated_by": userCtx.UserID,
	}

	if req.Name != "" {
		updateDoc["name"] = req.Name
	}
	if req.Description != "" {
		updateDoc["description"] = req.Description
	}
	if req.Image != "" {
		updateDoc["image"] = req.Image
	}
	if req.JobType != "" {
		updateDoc["job_type"] = req.JobType
	}
	if req.CronExpr != "" {
		updateDoc["cron_expr"] = req.CronExpr
	}
	if req.IntervalSeconds > 0 {
		updateDoc["interval_seconds"] = req.IntervalSeconds
	}
	if req.Command != "" {
		updateDoc["command"] = req.Command
	}
	if req.Status != "" {
		updateDoc["status"] = req.Status
	}
	if req.Timezone != "" {
		updateDoc["timezone"] = req.Timezone
	}

	// Increment generation for significant changes
	var updateQuery bson.M
	if req.JobType != "" || req.CronExpr != "" || req.IntervalSeconds > 0 {
		updateQuery = bson.M{
			"$set": updateDoc,
			"$inc": bson.M{"generation": 1},
		}
	} else {
		updateQuery = bson.M{"$set": updateDoc}
	}

	result, err := schedulersCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		updateQuery,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheduler"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scheduler not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scheduler updated successfully"})
}

func (h *SchedulerHandlers) DeleteScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(schedulerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduler ID"})
		return
	}

	schedulersCollection := h.db.Collection("scheduler_definitions")

	result, err := schedulersCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete scheduler"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scheduler not found"})
		return
	}

	// TODO: Also clean up related precompute and history records

	c.JSON(http.StatusOK, gin.H{"message": "Scheduler deleted successfully"})
}

func (h *SchedulerHandlers) RunScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	_, err := primitive.ObjectIDFromHex(schedulerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduler ID"})
		return
	}

	// Get current user from context
	_, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Implement immediate scheduler execution
	// This would typically involve:
	// 1. Creating a scheduler history record
	// 2. Dispatching the job to the execution engine
	// 3. Returning the run ID

	c.JSON(http.StatusOK, gin.H{
		"message": "Scheduler run initiated",
		"runId":   primitive.NewObjectID().Hex(),
	})
}

func (h *SchedulerHandlers) GetSchedulerHistory(c *gin.Context) {
	schedulerID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(schedulerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduler ID"})
		return
	}

	historyCollection := h.db.Collection("scheduler_history")

	cursor, err := historyCollection.Find(
		context.Background(),
		bson.M{"scheduler_id": objectID},
		options.Find().SetSort(bson.M{"start_time": -1}).SetLimit(50),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheduler history"})
		return
	}
	defer cursor.Close(context.Background())

	var history []models.SchedulerHistory
	if err := cursor.All(context.Background(), &history); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode scheduler history"})
		return
	}

	c.JSON(http.StatusOK, history)
}
