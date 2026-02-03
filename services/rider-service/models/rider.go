package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Rider ───────────────────────────────────────────────────
// One row per registered rider.  user_id is the FK that links
// back to auth_users (in auth-service's DB) – enforced at app level.

type Rider struct {
	ID            string         `gorm:"primaryKey;type:uuid"              json:"id"`
	UserID        string         `gorm:"uniqueIndex;not null;type:uuid"   json:"user_id"`
	Name          string         `gorm:"not null"                         json:"name"`
	Phone         string         `gorm:"index"                            json:"phone"`
	PaymentMethod string         `gorm:"default:'cash'"                   json:"payment_method"`
	TotalRides    int            `gorm:"default:0"                        json:"total_rides"`
	AvgRating     float64        `gorm:"default:0"                        json:"avg_rating"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"                   json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"                   json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"                            json:"deleted_at"`
}

func (r *Rider) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

func (Rider) TableName() string { return "riders" }
