package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ─── Config ──────────────────────────────────────────────────
// Each service passes its own env-derived Config to Connect().

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string // "disable" | "require"
}

// DSN builds the postgres connection string.
func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// ─── Connect ─────────────────────────────────────────────────
// Retries every 3 s until postgres is reachable (important in k8s
// where DB pods may start after app pods).

func Connect(cfg Config) *gorm.DB {
	customLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 30; i++ { // up to 90 s of retries
		db, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
			Logger: customLogger,
		})
		if err == nil {
			break
		}
		log.Printf("[DB] connection attempt %d failed: %v – retrying in 3s…", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("[DB] could not connect after 30 attempts: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("[DB] connected successfully")
	return db
}

// ─── AutoMigrate ─────────────────────────────────────────────
// Call once at startup with the service's model pointers.
// Example: db.AutoMigrate(User{}, Trip{})

func AutoMigrate(db *gorm.DB, models ...interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("[DB] auto-migrate failed: %v", err)
	}
	log.Println("[DB] auto-migrate complete")
}
