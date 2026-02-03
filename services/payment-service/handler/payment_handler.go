package handler

import (
	"context"
	"encoding/json"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/taxxy/payment-service/models"
	"github.com/taxxy/payment-service/repository"
	"github.com/taxxy/shared/messaging"
	paymentpb "github.com/taxxy/shared/proto/payment"
)

// ─── Server ──────────────────────────────────────────────────

type PaymentServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	repo repository.PaymentRepository
	mq   *messaging.Connection
}

func NewPaymentServer(repo repository.PaymentRepository, mq *messaging.Connection) *PaymentServer {
	return &PaymentServer{repo: repo, mq: mq}
}

// ─── CalculateFare ───────────────────────────────────────────

func (s *PaymentServer) CalculateFare(ctx context.Context, req *paymentpb.CalculateFareRequest) (*paymentpb.CalculateFareResponse, error) {
	distKm := haversine(req.PickupLat, req.PickupLng, req.DropoffLat, req.DropoffLng)

	// Estimate duration: assume avg speed 25 km/h in city
	durationMin := (distKm / 25.0) * 60.0

	// Tiered pricing
	rates := map[string]struct{ base, perKm, perMin float64 }{
		"sedan":  {8, 1.2, 0.3},
		"suv":    {12, 1.8, 0.4},
		"van":    {15, 2.2, 0.5},
		"luxury": {25, 3.0, 0.6},
	}
	r, ok := rates[req.VehicleType]
	if !ok {
		r = rates["sedan"]
	}

	baseFare    := r.base
	distFare    := distKm * r.perKm
	timeFare    := durationMin * r.perMin
	subtotal    := baseFare + distFare + timeFare
	surge       := 1.0 // TODO: dynamic surge based on demand
	totalFare   := math.Round(subtotal*surge*100) / 100

	breakdown := map[string]float64{
		"base_fare":   baseFare,
		"distance":    math.Round(distFare*100) / 100,
		"time":        math.Round(timeFare*100) / 100,
		"surge_mult":  surge,
		"total":       totalFare,
	}
	breakdownJSON, _ := json.Marshal(breakdown)

	return &paymentpb.CalculateFareResponse{
		BaseFare:        baseFare,
		DistanceKm:      math.Round(distKm*100) / 100,
		DurationMin:     math.Round(durationMin*10) / 10,
		SurgeMultiplier: surge,
		TotalFare:       totalFare,
		Breakdown:       string(breakdownJSON),
	}, nil
}

// ─── ProcessPayment ──────────────────────────────────────────

func (s *PaymentServer) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	if req.TripId == "" || req.RiderId == "" || req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "trip_id, rider_id and positive amount required")
	}

	payment := &models.Payment{
		TripID:        req.TripId,
		RiderID:       req.RiderId,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        models.PaymentPending,
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, status.Errorf(codes.Internal, "create payment: %v", err)
	}

	// Simulate processing (instant for card/wallet; pending for cash)
	finalStatus := models.PaymentCompleted
	if req.PaymentMethod == "cash" {
		finalStatus = models.PaymentPending // cash confirmed on completion
	}

	s.repo.UpdateStatus(payment.ID, finalStatus)

	// Publish payment.done event
	s.mq.Publish(messaging.EventPaymentDone, map[string]interface{}{
		"payment_id": payment.ID,
		"trip_id":    req.TripId,
		"rider_id":   req.RiderId,
		"amount":     req.Amount,
		"status":     finalStatus,
	})

	return &paymentpb.ProcessPaymentResponse{
		PaymentId: payment.ID,
		Status:    finalStatus,
		Amount:    req.Amount,
	}, nil
}

// ─── GetPayment ──────────────────────────────────────────────

func (s *PaymentServer) GetPayment(ctx context.Context, req *paymentpb.GetPaymentRequest) (*paymentpb.PaymentDetails, error) {
	p, err := s.repo.FindByID(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}
	return &paymentpb.PaymentDetails{
		PaymentId:     p.ID,
		TripId:        p.TripID,
		RiderId:       p.RiderID,
		Amount:        p.Amount,
		PaymentMethod: p.PaymentMethod,
		Status:        p.Status,
		CreatedAt:     p.CreatedAt.String(),
		UpdatedAt:     p.UpdatedAt.String(),
	}, nil
}

// ─── Refund ──────────────────────────────────────────────────

func (s *PaymentServer) Refund(ctx context.Context, req *paymentpb.RefundRequest) (*paymentpb.RefundResponse, error) {
	payment, err := s.repo.FindByID(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}
	if payment.Status != models.PaymentCompleted {
		return nil, status.Error(codes.FailedPrecondition, "can only refund completed payments")
	}
	if req.Amount > payment.Amount {
		return nil, status.Error(codes.InvalidArgument, "refund amount exceeds payment amount")
	}

	refund := &models.Refund{
		PaymentID: req.PaymentId,
		Amount:    req.Amount,
		Reason:    req.Reason,
		Status:    "completed", // instant refund in mock
	}

	if err := s.repo.CreateRefund(refund); err != nil {
		return nil, status.Errorf(codes.Internal, "create refund: %v", err)
	}

	// If full refund, mark payment as refunded
	if req.Amount == payment.Amount {
		s.repo.UpdateStatus(req.PaymentId, models.PaymentRefunded)
	}

	return &paymentpb.RefundResponse{
		RefundId: refund.ID,
		Status:   refund.Status,
		Amount:   refund.Amount,
	}, nil
}

// ─── GetDriverEarnings ───────────────────────────────────────

func (s *PaymentServer) GetDriverEarnings(ctx context.Context, req *paymentpb.GetDriverEarningsRequest) (*paymentpb.DriverEarnings, error) {
	total, trips, err := s.repo.SumEarningsByDriver(req.DriverId, req.FromDate, req.ToDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "earnings query: %v", err)
	}

	avg := 0.0
	if trips > 0 {
		avg = total / float64(trips)
	}

	return &paymentpb.DriverEarnings{
		DriverId:       req.DriverId,
		TotalEarned:    total,
		TotalTrips:     int32(trips),
		AveragePerTrip: math.Round(avg*100) / 100,
		PeriodFrom:     req.FromDate,
		PeriodTo:       req.ToDate,
	}, nil
}

// ─── haversine (same formula as trip-service; could extract to shared) ──

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371
	dLat := deg2rad(lat2 - lat1)
	dLng := deg2rad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func deg2rad(d float64) float64 { return d * math.Pi / 180 }
