package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI    string
	KafkaBroker string
	DBName      string
	Port        string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017/orchestrator"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		DBName:      getEnv("DB_NAME", "orchestrator"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
