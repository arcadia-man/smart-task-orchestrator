package api

import (
	"context"
	"net/http"
	"time"

	"smart-task-orchestrator/internal/models"
	"smart-task-orchestrator/internal/pkg/kafka"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobHandler struct {
	db       *mongo.Database
	producer *kafka.Producer
}

func NewJobHandler(db *mongo.Database, producer *kafka.Producer) *JobHandler {
	return &JobHandler{db: db, producer: producer}
}

type CreateJobRequest struct {
	Name     string                `json:"name" binding:"required"`
	Type     models.JobType        `json:"type" binding:"required"`
	Image    string                `json:"image" binding:"required"`
	Command  string                `json:"command" binding:"required"`
	CronExpr string                `json:"cron_expr"`
	Scaling  *models.SandboxConfig `json:"scaling"`
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var req CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	job := models.Job{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Image:     req.Image,
		Command:   req.Command,
		CronExpr:  req.CronExpr,
		Scaling:   req.Scaling,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := h.db.Collection("jobs").InsertOne(context.Background(), job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	// If it's a one-time job, trigger it immediately via Kafka
	if job.Type == models.JobTypeOneTime {
		err := h.producer.Publish(context.Background(), job.ID.Hex(), job)
		if err != nil {
			// Job created but failed to trigger - in real app, need a way to retry
			c.JSON(http.StatusAccepted, gin.H{"message": "Job created but failed to trigger immediately", "job": job})
			return
		}
	}

	c.JSON(http.StatusCreated, job)
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	cursor, err := h.db.Collection("jobs").Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}
	defer cursor.Close(context.Background())

	var jobs []models.Job
	if err := cursor.All(context.Background(), &jobs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode jobs"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, _ := primitive.ObjectIDFromHex(userIDStr)

	jobID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var job models.Job
	err = h.db.Collection("jobs").FindOne(context.Background(), bson.M{"_id": jobID, "user_id": userID}).Decode(&job)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

