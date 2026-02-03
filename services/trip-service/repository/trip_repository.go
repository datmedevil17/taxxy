package repository

import (
	"errors"

	"github.com/taxxy/trip-service/models"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("trip not found")

// ─── Interface ───────────────────────────────────────────────

type TripRepository interface {
	Create(trip *models.Trip) error
	FindByID(id string) (*models.Trip, error)
	UpdateStatus(tripID, status string) error
	AssignDriver(tripID, driverID string) error
	CompleteTripWithFare(tripID string, finalFare, distanceKm float64, durationSec int64) error
	ListByRider(riderID string, limit, offset int) ([]models.Trip, int64, error)
	ListByDriver(driverID string, limit, offset int) ([]models.Trip, int64, error)
	ListByStatus(status string, limit, offset int) ([]models.Trip, int64, error)
}

// ─── Impl ────────────────────────────────────────────────────

type tripRepo struct{ db *gorm.DB }

func NewTripRepository(db *gorm.DB) TripRepository {
	return &tripRepo{db: db}
}

func (r *tripRepo) Create(trip *models.Trip) error {
	return r.db.Create(trip).Error
}

func (r *tripRepo) FindByID(id string) (*models.Trip, error) {
	var trip models.Trip
	if err := r.db.Where("id = ?", id).First(&trip).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepo) UpdateStatus(tripID, status string) error {
	return r.db.Model(&models.Trip{}).
		Where("id = ?", tripID).
		Update("status", status).Error
}

func (r *tripRepo) AssignDriver(tripID, driverID string) error {
	return r.db.Model(&models.Trip{}).
		Where("id = ?", tripID).
		Updates(map[string]interface{}{
			"driver_id": driverID,
			"status":    models.StatusAccepted,
		}).Error
}

func (r *tripRepo) CompleteTripWithFare(tripID string, finalFare, distanceKm float64, durationSec int64) error {
	return r.db.Model(&models.Trip{}).
		Where("id = ?", tripID).
		Updates(map[string]interface{}{
			"status":       models.StatusCompleted,
			"final_fare":   finalFare,
			"distance_km":  distanceKm,
			"duration_sec": durationSec,
		}).Error
}

func (r *tripRepo) ListByRider(riderID string, limit, offset int) ([]models.Trip, int64, error) {
	var trips []models.Trip
	var count int64
	r.db.Model(&models.Trip{}).Where("rider_id = ?", riderID).Count(&count)
	err := r.db.Where("rider_id = ?", riderID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&trips).Error
	return trips, count, err
}

func (r *tripRepo) ListByDriver(driverID string, limit, offset int) ([]models.Trip, int64, error) {
	var trips []models.Trip
	var count int64
	r.db.Model(&models.Trip{}).Where("driver_id = ?", driverID).Count(&count)
	err := r.db.Where("driver_id = ?", driverID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&trips).Error
	return trips, count, err
}

func (r *tripRepo) ListByStatus(status string, limit, offset int) ([]models.Trip, int64, error) {
	var trips []models.Trip
	var count int64
	r.db.Model(&models.Trip{}).Where("status = ?", status).Count(&count)
	err := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&trips).Error
	return trips, count, err
}
