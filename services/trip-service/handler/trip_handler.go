package handler

import (
	"context"
	"fmt"
	"math"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/taxxy/shared/messaging"
	sharedpb "github.com/taxxy/shared/proto/shared"
	trippb "github.com/taxxy/shared/proto/trip"
	"github.com/taxxy/trip-service/models"
	"github.com/taxxy/trip-service/repository"
)

// ─── Server ──────────────────────────────────────────────────

type TripServer struct {
	trippb.UnimplementedTripServiceServer
	repo repository.TripRepository
	mq   *messaging.Connection
}

func NewTripServer(repo repository.TripRepository, mq *messaging.Connection) *TripServer {
	return &TripServer{repo: repo, mq: mq}
}

// ─── CreateTrip ──────────────────────────────────────────────

func (s *TripServer) CreateTrip(ctx context.Context, req *trippb.CreateTripRequest) (*trippb.CreateTripResponse, error) {
	if req.RiderId == "" || req.Pickup == nil || req.Dropoff == nil {
		return nil, status.Error(codes.InvalidArgument, "rider_id, pickup, dropoff required")
	}

	trip := &models.Trip{
		RiderID:       req.RiderId,
		PickupLat:     req.Pickup.Latitude,
		PickupLng:     req.Pickup.Longitude,
		DropoffLat:    req.Dropoff.Latitude,
		DropoffLng:    req.Dropoff.Longitude,
		VehicleType:   req.VehicleType,
		Status:        models.StatusRequested,
		EstimatedFare: req.EstimatedFare,
	}

	if err := s.repo.Create(trip); err != nil {
		return nil, status.Errorf(codes.Internal, "create trip: %v", err)
	}

	// ── Publish event so Driver Service can pick it up ──────
	s.mq.Publish(messaging.EventRideRequested, map[string]interface{}{
		"trip_id":        trip.ID,
		"rider_id":       trip.RiderID,
		"pickup_lat":     trip.PickupLat,
		"pickup_lng":     trip.PickupLng,
		"dropoff_lat":    trip.DropoffLat,
		"dropoff_lng":    trip.DropoffLng,
		"vehicle_type":   trip.VehicleType,
		"estimated_fare": trip.EstimatedFare,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	})

	return &trippb.CreateTripResponse{TripId: trip.ID, Status: trip.Status}, nil
}

// ─── GetTrip ─────────────────────────────────────────────────

func (s *TripServer) GetTrip(ctx context.Context, req *trippb.GetTripRequest) (*trippb.TripDetails, error) {
	trip, err := s.repo.FindByID(req.TripId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "trip not found")
	}
	return tripToProto(trip), nil
}

// ─── AssignDriver ────────────────────────────────────────────

func (s *TripServer) AssignDriver(ctx context.Context, req *trippb.AssignDriverRequest) (*trippb.AssignDriverResponse, error) {
	if req.TripId == "" || req.DriverId == "" {
		return nil, status.Error(codes.InvalidArgument, "trip_id and driver_id required")
	}

	// Guard: only "requested" trips can be assigned
	trip, err := s.repo.FindByID(req.TripId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "trip not found")
	}
	if trip.Status != models.StatusRequested {
		return nil, status.Errorf(codes.FailedPrecondition, "trip status is %s, expected requested", trip.Status)
	}

	if err := s.repo.AssignDriver(req.TripId, req.DriverId); err != nil {
		return nil, status.Errorf(codes.Internal, "assign driver: %v", err)
	}

	// ── Publish accepted event ──────────────────────────────
	s.mq.Publish(messaging.EventRideAccepted, map[string]interface{}{
		"trip_id":   req.TripId,
		"driver_id": req.DriverId,
		"rider_id":  trip.RiderID,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})

	return &trippb.AssignDriverResponse{TripId: req.TripId, Status: models.StatusAccepted}, nil
}

// ─── UpdateTripStatus ────────────────────────────────────────

func (s *TripServer) UpdateTripStatus(ctx context.Context, req *trippb.UpdateTripStatusRequest) (*trippb.UpdateTripStatusResponse, error) {
	trip, err := s.repo.FindByID(req.TripId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "trip not found")
	}

	// Simple state machine guard
	if !validTransition(trip.Status, req.Status) {
		return nil, status.Errorf(codes.FailedPrecondition,
			"cannot transition from %s to %s", trip.Status, req.Status)
	}

	var finalFare float64

	switch req.Status {
	case models.StatusCompleted:
		// Calculate actual fare based on haversine distance
		distKm := haversine(trip.PickupLat, trip.PickupLng, trip.DropoffLat, trip.DropoffLng)
		finalFare = calculateFare(distKm, trip.VehicleType)
		durSec := int64(time.Since(trip.CreatedAt).Seconds())

		if err := s.repo.CompleteTripWithFare(req.TripId, finalFare, distKm, durSec); err != nil {
			return nil, status.Errorf(codes.Internal, "complete trip: %v", err)
		}

		// ── Publish completed → triggers payment ────────────
		s.mq.Publish(messaging.EventRideCompleted, map[string]interface{}{
			"trip_id":    req.TripId,
			"rider_id":   trip.RiderID,
			"driver_id":  trip.DriverID,
			"final_fare": finalFare,
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		})

	case models.StatusCancelled:
		if err := s.repo.UpdateStatus(req.TripId, models.StatusCancelled); err != nil {
			return nil, status.Errorf(codes.Internal, "cancel trip: %v", err)
		}
		s.mq.Publish(messaging.EventRideCancelled, map[string]interface{}{
			"trip_id":   req.TripId,
			"rider_id":  trip.RiderID,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})

	default:
		if err := s.repo.UpdateStatus(req.TripId, req.Status); err != nil {
			return nil, status.Errorf(codes.Internal, "update status: %v", err)
		}
		if req.Status == models.StatusInProgress {
			s.mq.Publish(messaging.EventRideStarted, map[string]interface{}{
				"trip_id":   req.TripId,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		}
	}

	return &trippb.UpdateTripStatusResponse{
		TripId:    req.TripId,
		Status:    req.Status,
		FinalFare: finalFare,
	}, nil
}

// ─── ListTrips ───────────────────────────────────────────────

func (s *TripServer) ListTrips(ctx context.Context, req *trippb.ListTripsRequest) (*trippb.ListTripsResponse, error) {
	limit := int(req.Limit)
	offset := int(req.Offset)
	if limit == 0 {
		limit = 20
	}

	var trips []models.Trip
	var total int64
	var err error

	switch {
	case req.RiderId != "":
		trips, total, err = s.repo.ListByRider(req.RiderId, limit, offset)
	case req.DriverId != "":
		trips, total, err = s.repo.ListByDriver(req.DriverId, limit, offset)
	case req.Status != "":
		trips, total, err = s.repo.ListByStatus(req.Status, limit, offset)
	default:
		// No filter – return recent trips (admin / debug)
		trips, total, err = s.repo.ListByStatus("", limit, offset)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "list trips: %v", err)
	}

	var out []*trippb.TripDetails
	for i := range trips {
		out = append(out, tripToProto(&trips[i]))
	}

	return &trippb.ListTripsResponse{Trips: out, Total: int32(total)}, nil
}

// ─── helpers ─────────────────────────────────────────────────

func tripToProto(t *models.Trip) *trippb.TripDetails {
	return &trippb.TripDetails{
		TripId:        t.ID,
		RiderId:       t.RiderID,
		DriverId:      t.DriverID,
		Pickup:        &sharedpb.Location{Latitude: t.PickupLat, Longitude: t.PickupLng},
		Dropoff:       &sharedpb.Location{Latitude: t.DropoffLat, Longitude: t.DropoffLng},
		VehicleType:   t.VehicleType,
		Status:        t.Status,
		EstimatedFare: t.EstimatedFare,
		FinalFare:     t.FinalFare,
		DistanceKm:    t.DistanceKm,
		DurationSec:   t.DurationSec,
		CreatedAt:     t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     t.UpdatedAt.Format(time.RFC3339),
	}
}

// validTransition encodes the allowed state machine edges.
func validTransition(from, to string) bool {
	allowed := map[string][]string{
		models.StatusRequested:     {models.StatusFindingDriver, models.StatusCancelled},
		models.StatusFindingDriver: {models.StatusAccepted, models.StatusCancelled},
		models.StatusAccepted:      {models.StatusInProgress, models.StatusCancelled},
		models.StatusInProgress:    {models.StatusCompleted, models.StatusCancelled},
	}
	for _, s := range allowed[from] {
		if s == to {
			return true
		}
	}
	return false
}

// haversine returns the great-circle distance in km between two points.
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // earth radius km
	dLat := deg2rad(lat2 - lat1)
	dLng := deg2rad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func deg2rad(d float64) float64 { return d * math.Pi / 180 }

// calculateFare – simple tiered pricing per vehicle type.
// Base fare + per-km rate.  Surge can be layered later.
func calculateFare(distKm float64, vehicleType string) float64 {
	rates := map[string]struct{ base, perKm float64 }{
		"sedan":  {8, 1.5},
		"suv":    {12, 2.0},
		"van":    {15, 2.5},
		"luxury": {25, 3.5},
	}
	r, ok := rates[vehicleType]
	if !ok {
		r = rates["sedan"]
	}
	fare := r.base + distKm*r.perKm
	return math.Round(fare*100) / 100 // round to 2 dp
}

// Ensure fmt is used (for potential debug prints).
var _ = fmt.Sprintf
