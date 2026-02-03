package main

import (
	"encoding/json"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/taxxy/payment-service/config"
	"github.com/taxxy/payment-service/handler"
	"github.com/taxxy/payment-service/models"
	"github.com/taxxy/payment-service/repository"
	shareddb   "github.com/taxxy/shared/db"
	"github.com/taxxy/shared/messaging"
	paymentpb  "github.com/taxxy/shared/proto/payment"
)

func main() {
	cfg := config.Load()

	// ── Postgres ────────────────────────────────────────────
	db := shareddb.Connect(cfg.DB)
	shareddb.AutoMigrate(db, &models.Payment{}, &models.Refund{})

	// ── RabbitMQ ────────────────────────────────────────────
	mq := messaging.Connect(cfg.App.RabbitMQURL)
	defer mq.Close()

	repo := repository.NewPaymentRepository(db)

	// ── Subscribe to ride.completed → auto-process payment ─
	mq.Subscribe(
		"payment-service.ride.completed",
		[]string{messaging.EventRideCompleted},
		func(body []byte) error {
			var event struct {
				TripID    string  `json:"trip_id"`
				RiderID   string  `json:"rider_id"`
				DriverID  string  `json:"driver_id"`
				FinalFare float64 `json:"final_fare"`
			}
			if err := json.Unmarshal(body, &event); err != nil {
				return err
			}

			// Auto-create a payment record for completed trips
			payment := &models.Payment{
				TripID:        event.TripID,
				RiderID:       event.RiderID,
				Amount:        event.FinalFare,
				PaymentMethod: "card",  // default; would come from rider profile
				Status:        models.PaymentCompleted,
			}

			if err := repo.Create(payment); err != nil {
				return err
			}

			log.Printf("[PAYMENT-CONSUMER] auto-payment created: trip=%s amount=%.2f",
				event.TripID, event.FinalFare)
			return nil
		},
	)

	// ── Wire & serve gRPC ───────────────────────────────────
	paymentServer := handler.NewPaymentServer(repo, mq)

	lis, err := net.Listen("tcp", ":"+cfg.App.GRPCPort)
	if err != nil {
		log.Fatalf("[PAYMENT] listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(grpcServer, paymentServer)

	log.Printf("[PAYMENT] gRPC listening on :%s", cfg.App.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[PAYMENT] serve: %v", err)
	}
}
