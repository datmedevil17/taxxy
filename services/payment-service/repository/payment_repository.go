package repository

import (
	"errors"

	"github.com/taxxy/payment-service/models"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("payment not found")

// ─── Interface ───────────────────────────────────────────────

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByID(id string) (*models.Payment, error)
	FindByTripID(tripID string) (*models.Payment, error)
	UpdateStatus(paymentID, status string) error

	CreateRefund(refund *models.Refund) error
	SumEarningsByDriver(driverID, fromDate, toDate string) (float64, int64, error)
}

// ─── Impl ────────────────────────────────────────────────────

type paymentRepo struct{ db *gorm.DB }

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepo) FindByID(id string) (*models.Payment, error) {
	var p models.Payment
	if err := r.db.Where("id = ?", id).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) FindByTripID(tripID string) (*models.Payment, error) {
	var p models.Payment
	if err := r.db.Where("trip_id = ?", tripID).First(&p).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) UpdateStatus(paymentID, status string) error {
	return r.db.Model(&models.Payment{}).
		Where("id = ?", paymentID).
		Update("status", status).Error
}

func (r *paymentRepo) CreateRefund(refund *models.Refund) error {
	return r.db.Create(refund).Error
}

// SumEarningsByDriver – joins payments with trip_id and sums completed payments
// where the driver participated.  Note: driver_id is not on the payments table
// directly; in a real setup you'd join through trips or denormalize.
// For simplicity we accept a date range and query payments by date.
// The gateway will pass driver_id context; this is a simplified earnings calc.
func (r *paymentRepo) SumEarningsByDriver(driverID, fromDate, toDate string) (float64, int64, error) {
	// Simplified: in production you'd JOIN with trips table (cross-service)
	// or maintain a denormalized earnings table updated via the ride.completed event.
	// Here we return placeholder aggregation.
	var total float64
	var count int64

	// This would be a real query if earnings were denormalized here.
	// For now return zeros – the event consumer will populate a local earnings table.
	_ = driverID
	_ = fromDate
	_ = toDate

	return total, count, nil
}
