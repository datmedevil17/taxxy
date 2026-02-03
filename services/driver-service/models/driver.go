package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Driver ──────────────────────────────────────────────────

type Driver struct {
	ID           string         `gorm:"primaryKey;type:uuid"              json:"id"`
	UserID       string         `gorm:"uniqueIndex;not null;type:uuid"   json:"user_id"`
	Name         string         `gorm:"not null"                         json:"name"`
	Phone        string         `gorm:"index"                            json:"phone"`
	VehicleType  string         `gorm:"not null"                         json:"vehicle_type"`
	VehiclePlate string         `gorm:"not null"                         json:"vehicle_plate"`
	VehicleBrand string         `gorm:"not null"                         json:"vehicle_brand"`
	Status       string         `gorm:"not null;default:'offline'"       json:"status"` // available | busy | offline
	AvgRating    float64        `gorm:"default:0"                        json:"avg_rating"`
	TotalTrips   int            `gorm:"default:0"                        json:"total_trips"`
	// Current location (updated on availability ping)
	LastLat      float64        `gorm:"default:0"                        json:"last_lat"`
	LastLng      float64        `gorm:"default:0"                        json:"last_lng"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"                   json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"                   json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"                            json:"deleted_at"`
}

func (d *Driver) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return nil
}

func (Driver) TableName() string { return "drivers" }

// ─── DriverTrip ──────────────────────────────────────────────
// Local mirror: driver-service stores trips it participated in
// so we can answer GetMyTrips without always hitting trip-service.

type DriverTrip struct {
	ID        string    `gorm:"primaryKey;type:uuid"  json:"id"`
	DriverID  string    `gorm:"not null;index;type:uuid" json:"driver_id"`
	TripID    string    `gorm:"not null;uniqueIndex;type:uuid" json:"trip_id"`
	Status    string    `gorm:"not null;default:'accepted'" json:"status"`
	Fare      float64   `gorm:"default:0"            json:"fare"`
	CreatedAt time.Time `gorm:"autoCreateTime"       json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"       json:"updated_at"`
}

func (dt *DriverTrip) BeforeCreate(tx *gorm.DB) error {
	if dt.ID == "" {
		dt.ID = uuid.New().String()
	}
	return nil
}

func (DriverTrip) TableName() string { return "driver_trips" }
