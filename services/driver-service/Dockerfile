FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY shared/              ./shared/
COPY services/driver-service/ ./services/driver-service/
WORKDIR /app/services/driver-service
ENV GOTOOLCHAIN=auto
RUN go mod download && go build -o /driver-service .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /driver-service .
EXPOSE 50054
ENV GRPC_PORT=50054 \
    DB_HOST=postgres-driver \
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/ \
    TRIP_SERVICE_URL=trip-service:50053
ENTRYPOINT ["/app/driver-service"]
