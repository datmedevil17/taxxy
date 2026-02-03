package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/taxxy/trip-service/config"
	"github.com/taxxy/trip-service/handler"
	"github.com/taxxy/trip-service/models"
	"github.com/taxxy/trip-service/repository"
	shareddb  "github.com/taxxy/shared/db"
	"github.com/taxxy/shared/messaging"
	trippb    "github.com/taxxy/shared/proto/trip"
)

func main() {
	cfg := config.Load()

	// ── Postgres ────────────────────────────────────────────
	db := shareddb.Connect(cfg.DB)
	shareddb.AutoMigrate(db, &models.Trip{})

	// ── RabbitMQ ────────────────────────────────────────────
	mq := messaging.Connect(cfg.App.RabbitMQURL)
	defer mq.Close()

	// ── Wire & serve ────────────────────────────────────────
	repo       := repository.NewTripRepository(db)
	tripServer := handler.NewTripServer(repo, mq)

	lis, err := net.Listen("tcp", ":"+cfg.App.GRPCPort)
	if err != nil {
		log.Fatalf("[TRIP] listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	trippb.RegisterTripServiceServer(grpcServer, tripServer)

	log.Printf("[TRIP] gRPC listening on :%s", cfg.App.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[TRIP] serve: %v", err)
	}
}
