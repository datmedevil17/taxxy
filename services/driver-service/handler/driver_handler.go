package handler

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/taxxy/driver-service/config"
	"github.com/taxxy/driver-service/models"
	"github.com/taxxy/driver-service/repository"
	driverpb "github.com/taxxy/shared/proto/driver"
	trippb "github.com/taxxy/shared/proto/trip"
)

// ─── Server ──────────────────────────────────────────────────

type DriverServer struct {
	driverpb.UnimplementedDriverServiceServer
	repo repository.DriverRepository
	cfg  config.AppConfig
}

func NewDriverServer(repo repository.DriverRepository, cfg config.AppConfig) *DriverServer {
	return &DriverServer{repo: repo, cfg: cfg}
}

// ─── CreateProfile ───────────────────────────────────────────

func (s *DriverServer) CreateProfile(ctx context.Context, req *driverpb.CreateDriverProfileRequest) (*driverpb.CreateDriverProfileResponse, error) {
	if req.UserId == "" || req.Name == "" || req.VehicleType == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, name and vehicle_type required")
	}

	// Idempotent
	if existing, err := s.repo.FindByUserID(req.UserId); err == nil {
		return &driverpb.CreateDriverProfileResponse{DriverId: existing.ID, Status: "exists"}, nil
	}

	driver := &models.Driver{
		UserID:       req.UserId,
		Name:         req.Name,
		Phone:        req.Phone,
		VehicleType:  req.VehicleType,
		VehiclePlate: req.VehiclePlate,
		VehicleBrand: req.VehicleBrand,
		Status:       "offline",
	}

	if err := s.repo.Create(driver); err != nil {
		return nil, status.Errorf(codes.Internal, "create driver: %v", err)
	}

	return &driverpb.CreateDriverProfileResponse{DriverId: driver.ID, Status: "created"}, nil
}

// ─── GetProfile ──────────────────────────────────────────────

func (s *DriverServer) GetProfile(ctx context.Context, req *driverpb.GetDriverProfileRequest) (*driverpb.DriverProfile, error) {
	driver, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "driver not found")
	}
	return driverToProto(driver), nil
}

// ─── UpdateAvailability ──────────────────────────────────────

func (s *DriverServer) UpdateAvailability(ctx context.Context, req *driverpb.UpdateAvailabilityRequest) (*driverpb.UpdateAvailabilityResponse, error) {
	driver, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "driver not found")
	}

	if err := s.repo.UpdateStatus(driver.ID, req.Status); err != nil {
		return nil, status.Errorf(codes.Internal, "update status: %v", err)
	}

	// Also update location if provided
	if req.Location != nil {
		s.repo.UpdateLocation(driver.ID, req.Location.Latitude, req.Location.Longitude)
	}

	return &driverpb.UpdateAvailabilityResponse{Status: req.Status}, nil
}

// ─── AcceptRide ──────────────────────────────────────────────
// Driver explicitly accepts a pending trip.
// Calls Trip Service gRPC AssignDriver, then records locally.

func (s *DriverServer) AcceptRide(ctx context.Context, req *driverpb.AcceptRideRequest) (*driverpb.AcceptRideResponse, error) {
	driver, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "driver not found")
	}

	// ── Call Trip Service to assign this driver ──────────────
	tripClient, err := s.tripClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "trip service: %v", err)
	}

	_, err = tripClient.AssignDriver(ctx, &trippb.AssignDriverRequest{
		TripId:   req.TripId,
		DriverId: driver.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "assign driver on trip: %v", err)
	}

	// ── Record locally ───────────────────────────────────────
	s.repo.CreateDriverTrip(&models.DriverTrip{
		DriverID: driver.ID,
		TripID:   req.TripId,
		Status:   "accepted",
	})

	// Mark driver as busy
	s.repo.UpdateStatus(driver.ID, "busy")

	return &driverpb.AcceptRideResponse{TripId: req.TripId, Status: "accepted"}, nil
}

// ─── CompleteRide ────────────────────────────────────────────

func (s *DriverServer) CompleteRide(ctx context.Context, req *driverpb.CompleteRideRequest) (*driverpb.CompleteRideResponse, error) {
	driver, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "driver not found")
	}

	// Tell trip-service the trip is completed
	tripClient, err := s.tripClient()
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "trip service: %v", err)
	}

	resp, err := tripClient.UpdateTripStatus(ctx, &trippb.UpdateTripStatusRequest{
		TripId: req.TripId,
		Status: "completed",
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "complete trip: %v", err)
	}

	// Update local record
	s.repo.UpdateDriverTripStatus(req.TripId, "completed", resp.FinalFare)
	s.repo.IncrementTrips(driver.ID)
	s.repo.UpdateStatus(driver.ID, "available") // driver is free again

	return &driverpb.CompleteRideResponse{
		TripId: req.TripId,
		Status: "completed",
		Fare:   resp.FinalFare,
	}, nil
}

// ─── GetMyTrips ──────────────────────────────────────────────

func (s *DriverServer) GetMyTrips(ctx context.Context, req *driverpb.GetDriverTripsRequest) (*driverpb.GetDriverTripsResponse, error) {
	driver, err := s.repo.FindByUserID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "driver not found")
	}

	trips, total, err := s.repo.ListDriverTrips(driver.ID, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list trips: %v", err)
	}

	var out []*driverpb.DriverTripSummary
	for _, t := range trips {
		out = append(out, &driverpb.DriverTripSummary{
			TripId:    t.TripID,
			Status:    t.Status,
			Fare:      t.Fare,
			CreatedAt: t.CreatedAt.String(),
		})
	}

	return &driverpb.GetDriverTripsResponse{Trips: out, Total: int32(total)}, nil
}

// ─── helpers ─────────────────────────────────────────────────

func driverToProto(d *models.Driver) *driverpb.DriverProfile {
	return &driverpb.DriverProfile{
		DriverId:      d.ID,
		UserId:        d.UserID,
		Name:          d.Name,
		Phone:         d.Phone,
		VehicleType:   d.VehicleType,
		VehiclePlate:  d.VehiclePlate,
		VehicleBrand:  d.VehicleBrand,
		Status:        d.Status,
		AverageRating: d.AvgRating,
		TotalTrips:    int32(d.TotalTrips),
	}
}

func (s *DriverServer) tripClient() (trippb.TripServiceClient, error) {
	conn, err := grpc.NewClient(s.cfg.TripServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	log.Printf("[DRIVER] dialled trip-service at %s", s.cfg.TripServiceURL)
	return trippb.NewTripServiceClient(conn), nil
}
