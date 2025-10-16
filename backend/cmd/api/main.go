package main

import (
	"log"
	"net/http"

	"smart-task-orchestrator/internal/config"
	"smart-task-orchestrator/internal/db"
	"smart-task-orchestrator/internal/jobs"
	"smart-task-orchestrator/internal/kafka"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB
	mongoDB, err := db.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoDB.Close()

	// Initialize services
	jobService := jobs.NewService(mongoDB.Database)
	producer := kafka.NewProducer(cfg.KafkaBroker)
	defer producer.Close()

	// Setup Gin router
	r := gin.Default()

	// Simple CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Routes
	api := r.Group("/api")
	{
		api.POST("/jobs", createJob(jobService, producer))
		api.GET("/jobs", getAllJobs(jobService))
		api.GET("/jobs/:id", getJob(jobService))
		api.POST("/jobs/:id/retry", retryJob(jobService, producer))
		api.GET("/jobs/:id/status", getJobStatus(jobService))
	}

	log.Printf("🚀 API Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func createJob(jobService *jobs.Service, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req jobs.CreateJobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set default max retries
		if req.MaxRetries == 0 {
			req.MaxRetries = 3
		}

		job, err := jobService.CreateJob(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// For immediate jobs, publish to Kafka
		if req.Type == jobs.TypeImmediate {
			err = producer.PublishJob(c.Request.Context(), "jobs.execute", job.ID.Hex(), job.Payload)
			if err != nil {
				log.Printf("Failed to publish job to Kafka: %v", err)
			} else {
				jobService.UpdateJobStatus(c.Request.Context(), job.ID.Hex(), jobs.StatusQueued, "Job queued for execution")
			}
		}

		c.JSON(http.StatusCreated, job)
	}
}

func getAllJobs(jobService *jobs.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobs, err := jobService.GetAllJobs(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, jobs)
	}
}

func getJob(jobService *jobs.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		job, err := jobService.GetJobByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}

		c.JSON(http.StatusOK, job)
	}
}

func retryJob(jobService *jobs.Service, producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		job, err := jobService.GetJobByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}

		if job.RetryCount >= job.MaxRetries {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Job has exceeded maximum retries"})
			return
		}

		// Increment retry count
		err = jobService.IncrementRetryCount(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Publish to Kafka
		err = producer.PublishJob(c.Request.Context(), "jobs.execute", id, job.Payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue job for retry"})
			return
		}

		jobService.UpdateJobStatus(c.Request.Context(), id, jobs.StatusQueued, "Job queued for retry")
		c.JSON(http.StatusOK, gin.H{"message": "Job queued for retry"})
	}
}

func getJobStatus(jobService *jobs.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		job, err := jobService.GetJobByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":         job.ID,
			"status":     job.Status,
			"retryCount": job.RetryCount,
			"updatedAt":  job.UpdatedAt,
		})
	}
}
