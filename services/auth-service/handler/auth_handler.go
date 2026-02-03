package handler

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/taxxy/auth-service/config"
	"github.com/taxxy/auth-service/middleware"
	"github.com/taxxy/auth-service/models"
	"github.com/taxxy/auth-service/repository"
	authpb "github.com/taxxy/shared/proto/auth"
)

// ─── Server ──────────────────────────────────────────────────

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer // forward-compat
	repo                                  repository.UserRepository
	config                                config.AppConfig
}

func NewAuthServer(repo repository.UserRepository, cfg config.AppConfig) *AuthServer {
	return &AuthServer{repo: repo, config: cfg}
}

// ─── Register ────────────────────────────────────────────────

func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	// Validation
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password and name are required")
	}
	if req.Role != "rider" && req.Role != "driver" {
		return nil, status.Error(codes.InvalidArgument, "role must be 'rider' or 'driver'")
	}

	// Check uniqueness
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, status.Error(codes.AlreadyExists, "email already registered")
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password, // BeforeCreate hook will hash
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     req.Role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, status.Errorf(codes.Internal, "create user: %v", err)
	}

	token, err := middleware.GenerateToken(
		user.ID, user.Role, user.Email,
		s.config.JWTSecret,
		time.Duration(s.config.JWTExpiry)*time.Hour,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign token: %v", err)
	}

	return &authpb.RegisterResponse{UserId: user.ID, Token: token}, nil
}

// ─── Login ───────────────────────────────────────────────────

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	if err := user.ComparePassword(req.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	token, err := middleware.GenerateToken(
		user.ID, user.Role, user.Email,
		s.config.JWTSecret,
		time.Duration(s.config.JWTExpiry)*time.Hour,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sign token: %v", err)
	}

	return &authpb.LoginResponse{
		UserId: user.ID,
		Token:  token,
		Role:   user.Role,
	}, nil
}

// ─── ValidateToken ───────────────────────────────────────────
// Called by the API Gateway on every incoming request.

func (s *AuthServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	claims, err := middleware.ParseToken(req.Token, s.config.JWTSecret)
	if err != nil {
		return &authpb.ValidateTokenResponse{Valid: false}, nil // not an error – just invalid
	}
	return &authpb.ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Role:   claims.Role,
		Email:  claims.Email,
	}, nil
}

// ─── GetUser ─────────────────────────────────────────────────

func (s *AuthServer) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.GetUserResponse, error) {
	user, err := s.repo.FindByID(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}
	return &authpb.GetUserResponse{
		UserId: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Role:   user.Role,
		Phone:  user.Phone,
	}, nil
}
