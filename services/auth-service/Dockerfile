# ─── Stage 1: Build ──────────────────────────────────────────
FROM golang:1.23-alpine AS builder

# Install git (needed for some Go module fetches)
RUN apk add --no-cache git

WORKDIR /app

# Copy shared module first (layer-cache friendly)
COPY shared/go.mod    shared/go.sum    ./shared/
COPY shared/          ./shared/

# Copy service module files
COPY services/auth-service/go.mod  services/auth-service/go.sum  ./services/auth-service/
COPY services/auth-service/        ./services/auth-service/

# Download dependencies
WORKDIR /app/services/auth-service
ENV GOTOOLCHAIN=auto
RUN go mod download

# Build the binary
RUN go build -o /auth-service .

# ─── Stage 2: Runtime ────────────────────────────────────────
FROM alpine:3.19
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /auth-service .

# Expose the default gRPC port
EXPOSE 50051

# Default env (overridden by k8s ConfigMap / Secret)
ENV GRPC_PORT=50051 \
    DB_HOST=postgres-auth \
    DB_PORT=5432 \
    DB_USER=taxxy_auth \
    DB_NAME=taxxy_auth \
    DB_SSLMODE=disable \
    JWT_SECRET=taxxy-dev-secret-change-in-prod \
    JWT_EXPIRY_HOURS=24

ENTRYPOINT ["/app/auth-service"]
