package main

import (
	"encoding/json"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/taxxy/driver-service/config"
	"github.com/taxxy/driver-service/handler"
	"github.com/taxxy/driver-service/models"
	"github.com/taxxy/driver-service/repository"
	shareddb  "github.com/taxxy/shared/db"
	"github.com/taxxy/shared/messaging"
	driverpb  "github.com/taxxy/shared/proto/driver"
)

func main() {
	cfg := config.Load()

	// ── Postgres ────────────────────────────────────────────
	db := shareddb.Connect(cfg.DB)
	shareddb.AutoMigrate(db, &models.Driver{}, &models.DriverTrip{})

	// ── RabbitMQ ────────────────────────────────────────────
	mq := messaging.Connect(cfg.App.RabbitMQURL)
	defer mq.Close()

	// ── Subscribe to ride.requested ─────────────────────────
	// When a new ride comes in, driver-service logs it so the
	// driver app can poll / receive a push notification.
	// (In production you'd push via WebSocket; here we just log.)
	repo := repository.NewDriverRepository(db)

	mq.Subscribe(
		"driver-service.ride.requested",          // unique queue name
		[]string{messaging.EventRideRequested},   // routing keys
		func(body []byte) error {
			var event struct {
				TripID       string  `json:"trip_id"`
				VehicleType  string  `json:"vehicle_type"`
				PickupLat    float64 `json:"pickup_lat"`
				PickupLng    float64 `json:"pickup_lng"`
				EstimatedFare float64 `json:"estimated_fare"`
			}
			if err := json.Unmarshal(body, &event); err != nil {
				return err
			}

			// Find available drivers for this vehicle type
			drivers, err := repo.FindAvailableByVehicle(event.VehicleType)
			if err != nil {
				return err
			}

			log.Printf("[DRIVER-CONSUMER] ride.requested trip=%s vehicle=%s → %d available drivers",
				event.TripID, event.VehicleType, len(drivers))

			// In a real system you'd push a notification to each driver here.
			// For now this is the hook point for WebSocket / FCM integration.
			_ = drivers

			return nil
		},
	)

	// ── Wire & serve gRPC ───────────────────────────────────
	driverServer := handler.NewDriverServer(repo, cfg.App)

	lis, err := net.Listen("tcp", ":"+cfg.App.GRPCPort)
	if err != nil {
		log.Fatalf("[DRIVER] listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	driverpb.RegisterDriverServiceServer(grpcServer, driverServer)

	log.Printf("[DRIVER] gRPC listening on :%s", cfg.App.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[DRIVER] serve: %v", err)
	}
}
