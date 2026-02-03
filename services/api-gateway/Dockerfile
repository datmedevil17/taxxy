FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY shared/               ./shared/
COPY services/api-gateway/ ./services/api-gateway/
WORKDIR /app/services/api-gateway
RUN go mod download && go build -o /api-gateway .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /api-gateway .
EXPOSE 8080
ENV HTTP_PORT=8080 \
    AUTH_SERVICE_URL=auth-service:50051 \
    RIDER_SERVICE_URL=rider-service:50052 \
    TRIP_SERVICE_URL=trip-service:50053 \
    DRIVER_SERVICE_URL=driver-service:50054 \
    PAYMENT_SERVICE_URL=payment-service:50055
ENTRYPOINT ["/app/api-gateway"]
