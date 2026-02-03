package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Payment statuses ────────────────────────────────────────

const (
	PaymentPending   = "pending"
	PaymentCompleted = "completed"
	PaymentFailed    = "failed"
	PaymentRefunded  = "refunded"
)

// ─── Payment ─────────────────────────────────────────────────

type Payment struct {
	ID            string         `gorm:"primaryKey;type:uuid"       json:"id"`
	TripID        string         `gorm:"not null;index;type:uuid"   json:"trip_id"`
	RiderID       string         `gorm:"not null;index;type:uuid"   json:"rider_id"`
	Amount        float64        `gorm:"not null"                   json:"amount"`
	PaymentMethod string         `gorm:"not null"                   json:"payment_method"` // card|cash|wallet
	Status        string         `gorm:"not null;default:'pending'" json:"status"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"             json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"             json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"                      json:"deleted_at"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (Payment) TableName() string { return "payments" }

// ─── Refund ──────────────────────────────────────────────────

type Refund struct {
	ID        string    `gorm:"primaryKey;type:uuid"       json:"id"`
	PaymentID string    `gorm:"not null;index;type:uuid"   json:"payment_id"`
	Amount    float64   `gorm:"not null"                   json:"amount"`
	Reason    string    `gorm:"not null"                   json:"reason"`
	Status    string    `gorm:"not null;default:'pending'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime"             json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"             json:"updated_at"`
}

func (r *Refund) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

func (Refund) TableName() string { return "refunds" }
