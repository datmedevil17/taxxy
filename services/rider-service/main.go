package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/taxxy/rider-service/config"
	"github.com/taxxy/rider-service/handler"
	"github.com/taxxy/rider-service/models"
	"github.com/taxxy/rider-service/repository"
	shareddb "github.com/taxxy/shared/db"
	riderpb  "github.com/taxxy/shared/proto/rider"
)

func main() {
	cfg := config.Load()

	db := shareddb.Connect(cfg.DB)
	shareddb.AutoMigrate(db, &models.Rider{})

	repo        := repository.NewRiderRepository(db)
	riderServer := handler.NewRiderServer(repo, cfg.App)

	lis, err := net.Listen("tcp", ":"+cfg.App.GRPCPort)
	if err != nil {
		log.Fatalf("[RIDER] listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	riderpb.RegisterRiderServiceServer(grpcServer, riderServer)

	log.Printf("[RIDER] gRPC listening on :%s", cfg.App.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[RIDER] serve: %v", err)
	}
}
