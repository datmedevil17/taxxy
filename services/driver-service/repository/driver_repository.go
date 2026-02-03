package repository

import (
	"errors"

	"github.com/taxxy/driver-service/models"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("driver not found")

// ─── Interface ───────────────────────────────────────────────

type DriverRepository interface {
	Create(driver *models.Driver) error
	FindByUserID(userID string) (*models.Driver, error)
	FindByID(id string) (*models.Driver, error)
	UpdateStatus(driverID, status string) error
	UpdateLocation(driverID string, lat, lng float64) error
	FindAvailableByVehicle(vehicleType string) ([]models.Driver, error)
	IncrementTrips(driverID string) error

	// DriverTrip sub-repository
	CreateDriverTrip(dt *models.DriverTrip) error
	UpdateDriverTripStatus(tripID, status string, fare float64) error
	ListDriverTrips(driverID string, limit, offset int) ([]models.DriverTrip, int64, error)
}

// ─── Impl ────────────────────────────────────────────────────

type driverRepo struct{ db *gorm.DB }

func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverRepo{db: db}
}

func (r *driverRepo) Create(driver *models.Driver) error {
	return r.db.Create(driver).Error
}

func (r *driverRepo) FindByUserID(userID string) (*models.Driver, error) {
	var d models.Driver
	if err := r.db.Where("user_id = ?", userID).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *driverRepo) FindByID(id string) (*models.Driver, error) {
	var d models.Driver
	if err := r.db.Where("id = ?", id).First(&d).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &d, nil
}

func (r *driverRepo) UpdateStatus(driverID, status string) error {
	return r.db.Model(&models.Driver{}).
		Where("id = ?", driverID).
		Update("status", status).Error
}

func (r *driverRepo) UpdateLocation(driverID string, lat, lng float64) error {
	return r.db.Model(&models.Driver{}).
		Where("id = ?", driverID).
		Updates(map[string]interface{}{"last_lat": lat, "last_lng": lng}).Error
}

// FindAvailableByVehicle returns all online drivers matching the vehicle type.
// In production you'd add a geo-radius filter here (PostGIS or Redis Geo).
func (r *driverRepo) FindAvailableByVehicle(vehicleType string) ([]models.Driver, error) {
	var drivers []models.Driver
	err := r.db.Where("status = ? AND vehicle_type = ?", "available", vehicleType).
		Find(&drivers).Error
	return drivers, err
}

func (r *driverRepo) IncrementTrips(driverID string) error {
	return r.db.Model(&models.Driver{}).
		Where("id = ?", driverID).
		Update("total_trips", gorm.Expr("total_trips + 1")).Error
}

// ── DriverTrip helpers ──────────────────────────────────────

func (r *driverRepo) CreateDriverTrip(dt *models.DriverTrip) error {
	return r.db.Create(dt).Error
}

func (r *driverRepo) UpdateDriverTripStatus(tripID, status string, fare float64) error {
	return r.db.Model(&models.DriverTrip{}).
		Where("trip_id = ?", tripID).
		Updates(map[string]interface{}{"status": status, "fare": fare}).Error
}

func (r *driverRepo) ListDriverTrips(driverID string, limit, offset int) ([]models.DriverTrip, int64, error) {
	var trips []models.DriverTrip
	var count int64
	r.db.Model(&models.DriverTrip{}).Where("driver_id = ?", driverID).Count(&count)
	err := r.db.Where("driver_id = ?", driverID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&trips).Error
	return trips, count, err
}
