package config

import (
	"os"
)

type Config struct {
	// Database
	MongoURI string
	DBName   string

	// Cache & Queue
	RedisURL    string
	KafkaBroker string

	// Authentication
	JWTSecret string

	// Server
	Port string

	// Docker
	DockerHost string
}

func Load() *Config {
	return &Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017/orchestrator"),
		DBName:      getEnv("DB_NAME", "orchestrator"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Port:        getEnv("PORT", "8080"),
		DockerHost:  getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
