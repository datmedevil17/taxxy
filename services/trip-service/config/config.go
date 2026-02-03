package config

import (
	"os"

	shareddb "github.com/taxxy/shared/db"
)

type AppConfig struct {
	GRPCPort     string
	RabbitMQURL  string
}

type Config struct {
	App AppConfig
	DB  shareddb.Config
}

func Load() Config {
	return Config{
		App: AppConfig{
			GRPCPort:    getEnv("GRPC_PORT", "50053"),
			RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		},
		DB: shareddb.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "taxxy_trip"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "taxxy_trip"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
