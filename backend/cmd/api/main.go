package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"smart-task-orchestrator/internal/auth"
	"smart-task-orchestrator/internal/config"
	"smart-task-orchestrator/internal/db"
	"smart-task-orchestrator/internal/handlers"
	"smart-task-orchestrator/pkg/kafka"
	"smart-task-orchestrator/pkg/redis"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	mongoDB, err := db.NewMongoDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize Redis
	redisClient, err := redis.NewClient(cfg.RedisURL, 4) // 4 shards
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize Kafka producer
	producer := kafka.NewProducer(cfg.KafkaBroker, "job_executions")
	defer producer.Close()

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret)

	// Initialize handlers
	authHandlers := handlers.NewAuthHandlers(mongoDB.Database, jwtManager)
	userHandlers := handlers.NewUserHandlers(mongoDB.Database)
	roleHandlers := handlers.NewRoleHandlers(mongoDB.Database)
	schedulerHandlers := handlers.NewSchedulerHandlers(mongoDB.Database)
	imageHandlers := handlers.NewImageHandlers(mongoDB.Database)
	logHandlers := handlers.NewLogHandlers(mongoDB.Database)
	monitoringHandlers := handlers.NewMonitoringHandlers()

	// Setup Gin router
	if gin.Mode() == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"service":   "api",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Authentication routes (no auth required)
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/login", authHandlers.Login)
			authRoutes.POST("/refresh", authHandlers.RefreshToken)
			authRoutes.POST("/change-password", auth.AuthMiddleware(jwtManager), authHandlers.ChangePassword)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware(jwtManager))
		{
			// User info
			protected.GET("/me", authHandlers.Me)

			// Dashboard stats
			protected.GET("/dashboard/stats", schedulerHandlers.GetSchedulerStats)

			// Schedulers
			schedulers := protected.Group("/schedulers")
			{
				schedulers.GET("", schedulerHandlers.GetSchedulers)
				schedulers.POST("", schedulerHandlers.CreateScheduler)
				schedulers.GET("/:id", schedulerHandlers.GetScheduler)
				schedulers.PUT("/:id", schedulerHandlers.UpdateScheduler)
				schedulers.DELETE("/:id", schedulerHandlers.DeleteScheduler)
				schedulers.POST("/:id/run", schedulerHandlers.RunScheduler)
				schedulers.GET("/:id/history", schedulerHandlers.GetSchedulerHistory)
			}

			// Users (admin only)
			users := protected.Group("/users")
			{
				users.GET("", userHandlers.GetUsers)
				users.POST("", userHandlers.CreateUser)
				users.PUT("/:id", userHandlers.UpdateUser)
				users.DELETE("/:id", userHandlers.DeleteUser)
				users.POST("/:id/reset-password", userHandlers.ResetPassword)
			}

			// Roles
			roles := protected.Group("/roles")
			{
				roles.GET("", roleHandlers.GetRoles)
				roles.POST("", roleHandlers.CreateRole)
				roles.PUT("/:id", roleHandlers.UpdateRole)
				roles.DELETE("/:id", roleHandlers.DeleteRole)
				roles.GET("/permissions", roleHandlers.GetPermissions)
			}

			// Images
			images := protected.Group("/images")
			{
				images.GET("", imageHandlers.GetImages)
				images.POST("", imageHandlers.CreateImage)
				images.PUT("/:id", imageHandlers.UpdateImage)
				images.DELETE("/:id", imageHandlers.DeleteImage)
			}

			// Logs
			logs := protected.Group("/logs")
			{
				logs.GET("", logHandlers.GetLogs)
				logs.GET("/stats", logHandlers.GetLogStats)
				logs.GET("/sources", logHandlers.GetLogSources)
			}

			// Monitoring
			monitoring := protected.Group("/monitoring")
			{
				monitoring.GET("", monitoringHandlers.GetFullMonitoring)
				monitoring.GET("/metrics", monitoringHandlers.GetSystemMetrics)
				monitoring.GET("/services", monitoringHandlers.GetServices)
				monitoring.GET("/alerts", monitoringHandlers.GetAlerts)
			}
		}
	}

	// WebSocket routes for real-time logs
	router.GET("/ws/logs/:runId", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "WebSocket logs - to be implemented"}) })

	// Start server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("🛑 Shutting down API server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}
	}()

	log.Printf("🚀 API Server starting on port %s", cfg.Port)
	log.Printf("📱 Frontend: http://localhost:3000")
	log.Printf("🔧 API: http://localhost:%s", cfg.Port)
	log.Printf("📊 Health: http://localhost:%s/health", cfg.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Failed to start server: %v", err)
	}

	log.Println("✅ API server stopped")
}
