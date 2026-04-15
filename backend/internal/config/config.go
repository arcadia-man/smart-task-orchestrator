package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DBUri        string `mapstructure:"DB_URI"`
	DBName       string `mapstructure:"DB_NAME"`
	KafkaBrokers string `mapstructure:"KAFKA_BROKERS"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	Environment  string `mapstructure:"ENVIRONMENT"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_URI", "mongodb://localhost:27017")
	viper.SetDefault("DB_NAME", "smart_orchestrator")
	viper.SetDefault("KAFKA_BROKERS", "localhost:9092")
	viper.SetDefault("JWT_SECRET", "super-secret-key")
	viper.SetDefault("ENVIRONMENT", "development")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
