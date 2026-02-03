package repository

import (
	"errors"

	"github.com/taxxy/rider-service/models"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("rider not found")

// ─── Interface ───────────────────────────────────────────────

type RiderRepository interface {
	Create(rider *models.Rider) error
	FindByUserID(userID string) (*models.Rider, error)
	FindByID(id string) (*models.Rider, error)
	IncrementRides(riderID string) error
}

// ─── Impl ────────────────────────────────────────────────────

type riderRepo struct{ db *gorm.DB }

func NewRiderRepository(db *gorm.DB) RiderRepository {
	return &riderRepo{db: db}
}

func (r *riderRepo) Create(rider *models.Rider) error {
	return r.db.Create(rider).Error
}

func (r *riderRepo) FindByUserID(userID string) (*models.Rider, error) {
	var rider models.Rider
	if err := r.db.Where("user_id = ?", userID).First(&rider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rider, nil
}

func (r *riderRepo) FindByID(id string) (*models.Rider, error) {
	var rider models.Rider
	if err := r.db.Where("id = ?", id).First(&rider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &rider, nil
}

func (r *riderRepo) IncrementRides(riderID string) error {
	return r.db.Model(&models.Rider{}).
		Where("id = ?", riderID).
		Update("total_rides", gorm.Expr("total_rides + 1")).
		Error
}
