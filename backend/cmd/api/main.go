package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"smart-task-orchestrator/internal/api"
	"smart-task-orchestrator/internal/config"
	"smart-task-orchestrator/internal/pkg/db"
	"smart-task-orchestrator/internal/pkg/kafka"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	client, err := db.ConnectMongo(cfg.DBUri)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	database := client.Database(cfg.DBName)

	// Kafka producer
	producer := kafka.NewProducer(strings.Split(cfg.KafkaBrokers, ","), "jobs.execute")
	defer producer.Close()

	r := gin.Default()

	// CORS - allow frontend dev server
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check (used by docker-compose)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes
	authHandler := api.NewAuthHandler(database, cfg.JWTSecret)
	r.POST("/api/signup", authHandler.Signup)
	r.POST("/api/login", authHandler.Login)

	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(api.AuthMiddleware(cfg.JWTSecret, database))
	{
		jobHandler := api.NewJobHandler(database, producer)
		authorized.POST("/jobs", jobHandler.CreateJob)
		authorized.GET("/jobs", jobHandler.ListJobs)
		authorized.GET("/jobs/:id", jobHandler.GetJob)

		authorized.GET("/profile", authHandler.GetProfile)
		authorized.PUT("/profile", authHandler.UpdateProfile)
		authorized.POST("/reset-password", authHandler.ResetPassword)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

