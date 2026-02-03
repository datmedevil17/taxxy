package config

import (
	"os"

	shareddb "github.com/taxxy/shared/db"
)

type AppConfig struct {
	GRPCPort      string
	TripServiceURL string // grpc://trip-service:50053
}

type Config struct {
	App AppConfig
	DB  shareddb.Config
}

func Load() Config {
	return Config{
		App: AppConfig{
			GRPCPort:       getEnv("GRPC_PORT", "50052"),
			TripServiceURL: getEnv("TRIP_SERVICE_URL", "trip-service:50053"),
		},
		DB: shareddb.Config{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "taxxy_rider"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "taxxy_rider"),
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
