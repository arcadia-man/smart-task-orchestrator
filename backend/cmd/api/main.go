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
    "go.mongodb.org/mongo-driver/mongo"

    "smart-task-orchestrator/internal/auth"
    "smart-task-orchestrator/internal/config"
    "smart-task-orchestrator/internal/db"
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
            authRoutes.POST("/login", loginHandler(mongoDB.Database, jwtManager))
            authRoutes.POST("/refresh", refreshHandler(jwtManager))
            authRoutes.POST("/change-password", auth.AuthMiddleware(jwtManager), changePasswordHandler(mongoDB.Database))
        }

        // Protected routes
        protected := api.Group("/")
        protected.Use(auth.AuthMiddleware(jwtManager))
        {
            // User info
            protected.GET("/me", meHandler(mongoDB.Database))

            // Schedulers
            schedulers := protected.Group("/schedulers")
            {
                schedulers.GET("", getSchedulersHandler(mongoDB.Database))
                schedulers.POST("", createSchedulerHandler(mongoDB.Database))
                schedulers.GET("/:id", getSchedulerHandler(mongoDB.Database))
                schedulers.PUT("/:id", updateSchedulerHandler(mongoDB.Database))
                schedulers.DELETE("/:id", deleteSchedulerHandler(mongoDB.Database))
                schedulers.POST("/:id/run", runSchedulerHandler(mongoDB.Database, producer))
                schedulers.GET("/:id/history", getSchedulerHistoryHandler(mongoDB.Database))
            }

            // Users (admin only)
            users := protected.Group("/users")
            {
                users.GET("", getUsersHandler(mongoDB.Database))
                users.POST("", createUserHandler(mongoDB.Database))
                users.PUT("/:id", updateUserHandler(mongoDB.Database))
                users.DELETE("/:id", deleteUserHandler(mongoDB.Database))
            }

            // Roles
            roles := protected.Group("/roles")
            {
                roles.GET("", getRolesHandler(mongoDB.Database))
                roles.POST("", createRoleHandler(mongoDB.Database))
            }

            // Images
            images := protected.Group("/images")
            {
                images.GET("", getImagesHandler(mongoDB.Database))
                images.POST("", createImageHandler(mongoDB.Database))
            }
        }
    }

    // WebSocket routes for real-time logs
    router.GET("/ws/logs/:runId", websocketHandler(redisClient))

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

// Placeholder handlers - these will be implemented in separate files
func loginHandler(db *mongo.Database, jwtManager *auth.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Login handler - to be implemented"})
    }
}

func refreshHandler(jwtManager *auth.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Refresh handler - to be implemented"})
    }
}

func changePasswordHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Change password handler - to be implemented"})
    }
}

func meHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Me handler - to be implemented"})
    }
}

func getSchedulersHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Get schedulers handler - to be implemented"})
    }
}

func createSchedulerHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Create scheduler handler - to be implemented"})
    }
}

func getSchedulerHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Get scheduler handler - to be implemented"})
    }
}

func updateSchedulerHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Update scheduler handler - to be implemented"})
    }
}

func deleteSchedulerHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Delete scheduler handler - to be implemented"})
    }
}

func runSchedulerHandler(db *mongo.Database, producer *kafka.Producer) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Run scheduler handler - to be implemented"})
    }
}

func getSchedulerHistoryHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Get scheduler history handler - to be implemented"})
    }
}

func getUsersHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Get users handler - to be implemented"})
    }
}

func createUserHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Create user handler - to be implemented"})
    }
}

func updateUserHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Update user handler - to be implemented"})
    }
}

func deleteUserHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Delete user handler - to be implemented"})
    }
}

func getRolesHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Get roles handler - to be implemented"})
    }
}

func createRoleHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Create role handler - to be implemented"})
    }
}

func getImagesHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Get images handler - to be implemented"})
    }
}

func createImageHandler(db *mongo.Database) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Create image handler - to be implemented"})
    }
}

func websocketHandler(redisClient *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "WebSocket handler - to be implemented"})
    }
}
