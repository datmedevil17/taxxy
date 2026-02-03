package handler

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/taxxy/rider-service/config"
	"github.com/taxxy/rider-service/models"
	"github.com/taxxy/rider-service/repository"
	riderpb "github.com/taxxy/shared/proto/rider"
	sharedpb "github.com/taxxy/shared/proto/shared"
	trippb "github.com/taxxy/shared/proto/trip"
)

// ─── Server ──────────────────────────────────────────────────

type RiderServer struct {
	riderpb.UnimplementedRiderServiceServer
	repo repository.RiderRepository
	cfg  config.AppConfig
}

func NewRiderServer(repo repository.RiderRepository, cfg config.AppConfig) *RiderServer {
	return &RiderServer{repo: repo, cfg: cfg}
}

// ─── CreateProfile ───────────────────────────────────────────

func (s *RiderServer) CreateProfile(ctx context.Context, req *riderpb.CreateRiderProfileRequest) (*riderpb.CreateRiderProfileResponse, error) {
	if req.UserId == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and name are required")
	}

	// Idempotency: if profile already exists, return it.
	if existing, err := s.repo.FindByUserID(req.UserId); err == nil {
		return &riderpb.CreateRiderProfileResponse{RiderId: existing.ID, Status: "exists"}, nil
	}

	rider := &models.Rider{
		UserID:        req.UserId,
		Name:          req.Name,
		Phone:         req.Phone,
		PaymentMethod: req.PaymentMethod,
	}

	if err := s.repo.Create(rider); err != nil {
		return nil, status.Errorf(codes.Internal, "create rider: %v", err)
	}

	return &riderpb.CreateRiderProfileResponse{RiderId: rider.ID, Status: "created"}, nil
}

// ─── GetProfile ──────────────────────────────────────────────

func (s *RiderServer) GetProfile(ctx context.Context, req *riderpb.GetRiderProfileRequest) (*riderpb.RiderProfile, error) {
	rider, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "rider profile not found")
	}
	return &riderpb.RiderProfile{
		RiderId:       rider.ID,
		UserId:        rider.UserID,
		Name:          rider.Name,
		Phone:         rider.Phone,
		PaymentMethod: rider.PaymentMethod,
		TotalRides:    int32(rider.TotalRides),
		AverageRating: rider.AvgRating,
	}, nil
}

// ─── RequestRide ─────────────────────────────────────────────
// Orchestrates: validate rider → call Trip Service → return trip ID.

func (s *RiderServer) RequestRide(ctx context.Context, req *riderpb.RequestRideRequest) (*riderpb.RequestRideResponse, error) {
	// 1. Ensure rider profile exists
	rider, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "rider profile not found – create profile first")
	}

	// 2. Validate locations
	if req.Pickup == nil || req.Dropoff == nil {
		return nil, status.Error(codes.InvalidArgument, "pickup and dropoff locations required")
	}

	// 3. Call Trip Service to create the trip
	tripClient, err := s.tripClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "trip service: %v", err)
	}

	// Estimate fare (simple mock: $2/km base; real calc in payment-service)
	estimatedFare := 15.0 // placeholder – payment-service.CalculateFare will refine

	tripResp, err := tripClient.CreateTrip(ctx, &trippb.CreateTripRequest{
		RiderId:       rider.ID,
		Pickup:        &sharedpb.Location{Latitude: req.Pickup.Latitude, Longitude: req.Pickup.Longitude},
		Dropoff:       &sharedpb.Location{Latitude: req.Dropoff.Latitude, Longitude: req.Dropoff.Longitude},
		VehicleType:   req.VehicleType,
		EstimatedFare: estimatedFare,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create trip: %v", err)
	}

	return &riderpb.RequestRideResponse{
		TripId:        tripResp.TripId,
		Status:        tripResp.Status,
		EstimatedFare: estimatedFare,
		Eta:           "5 min", // placeholder
	}, nil
}

// ─── GetMyTrips ──────────────────────────────────────────────

func (s *RiderServer) GetMyTrips(ctx context.Context, req *riderpb.GetRiderTripsRequest) (*riderpb.GetRiderTripsResponse, error) {
	rider, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "rider not found")
	}

	// Call Trip Service to list trips for this rider
	tripClient, err := s.tripClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "trip service: %v", err)
	}

	listResp, err := tripClient.ListTrips(ctx, &trippb.ListTripsRequest{
		RiderId: rider.ID,
		Limit:   req.Limit,
		Offset:  req.Offset,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list trips: %v", err)
	}

	var summaries []*riderpb.RiderTripSummary
	for _, t := range listResp.Trips {
		summaries = append(summaries, &riderpb.RiderTripSummary{
			TripId:    t.TripId,
			Status:    t.Status,
			Fare:      t.FinalFare,
			CreatedAt: t.CreatedAt,
		})
	}

	return &riderpb.GetRiderTripsResponse{Trips: summaries, Total: listResp.Total}, nil
}

// ─── helpers ─────────────────────────────────────────────────

func (s *RiderServer) tripClient() (trippb.TripServiceClient, error) {
	conn, err := grpc.NewClient(s.cfg.TripServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	// Note: in production you'd cache this connection. Fine for clarity here.
	log.Printf("[RIDER] dialled trip-service at %s", s.cfg.TripServiceURL)
	return trippb.NewTripServiceClient(conn), nil
}
