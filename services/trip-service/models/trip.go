package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Status constants ─────────────────────────────────────────

const (
	StatusRequested     = "requested"
	StatusFindingDriver = "finding_driver"
	StatusAccepted      = "accepted"
	StatusInProgress    = "in_progress"
	StatusCompleted     = "completed"
	StatusCancelled     = "cancelled"
)

// ─── Trip ─────────────────────────────────────────────────────
// Central record for a single ride.  Pickup/Dropoff are stored
// as flat lat/lng columns (no PostGIS dependency for v1).

type Trip struct {
	ID            string         `gorm:"primaryKey;type:uuid"       json:"id"`
	RiderID       string         `gorm:"not null;index;type:uuid"   json:"rider_id"`
	DriverID      *string        `gorm:"index;type:uuid"            json:"driver_id"` // NULL until accepted
	PickupLat     float64        `gorm:"not null"                   json:"pickup_lat"`
	PickupLng     float64        `gorm:"not null"                   json:"pickup_lng"`
	DropoffLat    float64        `gorm:"not null"                   json:"dropoff_lat"`
	DropoffLng    float64        `gorm:"not null"                   json:"dropoff_lng"`
	VehicleType   string         `gorm:"not null"                   json:"vehicle_type"`
	Status        string         `gorm:"not null;default:requested" json:"status"`
	EstimatedFare float64        `gorm:"default:0"                 json:"estimated_fare"`
	FinalFare     float64        `gorm:"default:0"                 json:"final_fare"`
	DistanceKm    float64        `gorm:"default:0"                 json:"distance_km"`
	DurationSec   int64          `gorm:"default:0"                 json:"duration_sec"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"            json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"            json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"                     json:"deleted_at"`
}

func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

func (Trip) TableName() string { return "trips" }
