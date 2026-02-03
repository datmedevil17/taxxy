.PHONY: all build run stop test clean proto frontend install help

# Default target
all: help

# ─── Help ─────────────────────────────────────────────────────────────
help:
	@echo "Available commands:"
	@echo "  make proto      - Generate Go code from Protocol Buffers"
	@echo "  make build      - Build all backend services"
	@echo "  make up         - Start all services with Docker Compose (builds if needed)"
	@echo "  make down       - Stop all services"
	@echo "  make logs       - View logs from all services"
	@echo "  make frontend   - Run the frontend development server"
	@echo "  make install    - Install frontend dependencies"
	@echo "  make clean      - Clean up artifacts and temporary files"
	@echo ""
	@echo "Database access:"
	@echo "  make db-auth    - Connect to Auth database"
	@echo "  make db-rider   - Connect to Rider database"
	@echo "  make db-trip    - Connect to Trip database"
	@echo "  make db-driver  - Connect to Driver database"
	@echo "  make db-payment - Connect to Payment database"

# ─── Proto ────────────────────────────────────────────────────────────
proto:
	@echo "Generating proto code..."
	@protoc proto/shared/shared.proto --go_out=. --go-grpc_out=.
	@protoc proto/driver/driver.proto --go_out=. --go-grpc_out=.
	@protoc proto/trip/trip.proto --go_out=. --go-grpc_out=.
	@protoc proto/rider/rider.proto --go_out=. --go-grpc_out=.
	@protoc proto/payment/payment.proto --go_out=. --go-grpc_out=.
	@echo "Proto generation complete."

# ─── Docker ───────────────────────────────────────────────────────────
up:
	@echo "Starting services..."
	@docker compose up --build -d
	@echo "Services started. API Gateway at http://localhost:8080"

down:
	@echo "Stopping services..."
	@docker compose down

logs:
	@docker compose logs -f

# ─── Frontend ─────────────────────────────────────────────────────────
install:
	@echo "Installing frontend dependencies..."
	@cd web && npm install

frontend:
	@echo "Starting frontend dev server..."
	@cd web && npm run dev

# ─── Go ───────────────────────────────────────────────────────────────
build:
	@echo "Building services..."
	@go build -o bin/api-gateway ./services/api-gateway
	@go build -o bin/auth-service ./services/auth-service
	@go build -o bin/driver-service ./services/driver-service
	@go build -o bin/trip-service ./services/trip-service
	@go build -o bin/payment-service ./services/payment-service
	@go build -o bin/rider-service ./services/rider-service
	@echo "Build complete."

clean:
	@echo "Cleaning up..."
	@rm -rf bin
	@rm -rf web/dist

# ─── Database Access ──────────────────────────────────────────────────
db-auth:
	@docker exec -it taxxy-postgres-auth-1 psql -U taxxy_auth -d taxxy_auth

db-rider:
	@docker exec -it taxxy-postgres-rider-1 psql -U taxxy_rider -d taxxy_rider

db-trip:
	@docker exec -it taxxy-postgres-trip-1 psql -U taxxy_trip -d taxxy_trip

db-driver:
	@docker exec -it taxxy-postgres-driver-1 psql -U taxxy_driver -d taxxy_driver

db-payment:
	@docker exec -it taxxy-postgres-payment-1 psql -U taxxy_payment -d taxxy_payment
