package config

import (
	"os"
	"strconv"

	shareddb "github.com/taxxy/shared/db"
)

// ─── App ─────────────────────────────────────────────────────

type AppConfig struct {
	GRPCPort   string
	JWTSecret  string
	JWTExpiry  int // hours
}

// ─── DB ──────────────────────────────────────────────────────

type DBConfig = shareddb.Config

// ─── Root ────────────────────────────────────────────────────

type Config struct {
	App AppConfig
	DB  DBConfig
}

func Load() Config {
	return Config{
		App: AppConfig{
			GRPCPort:  getEnv("GRPC_PORT", "50051"),
			JWTSecret: getEnv("JWT_SECRET", "taxxy-dev-secret-change-in-prod"),
			JWTExpiry: getEnvInt("JWT_EXPIRY_HOURS", 24),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "taxxy_auth"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "taxxy_auth"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// ─── helpers ─────────────────────────────────────────────────

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
