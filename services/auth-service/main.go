package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/taxxy/auth-service/config"
	"github.com/taxxy/auth-service/handler"
	"github.com/taxxy/auth-service/models"
	"github.com/taxxy/auth-service/repository"
	shareddb "github.com/taxxy/shared/db"
	authpb "github.com/taxxy/shared/proto/auth"
)

func main() {
	// ── 1. Load config ──────────────────────────────────────
	cfg := config.Load()

	// ── 2. Connect to Postgres & auto-migrate ───────────────
	db := shareddb.Connect(cfg.DB)
	shareddb.AutoMigrate(db, &models.User{})

	// ── 3. Wire dependencies ────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	authServer := handler.NewAuthServer(userRepo, cfg.App)

	// ── 4. Start gRPC server ────────────────────────────────
	lis, err := net.Listen("tcp", ":"+cfg.App.GRPCPort)
	if err != nil {
		log.Fatalf("[AUTH] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, authServer)

	log.Printf("[AUTH] gRPC server listening on :%s", cfg.App.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[AUTH] gRPC serve failed: %v", err)
	}
}
