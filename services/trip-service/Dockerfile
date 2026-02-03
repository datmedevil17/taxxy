FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY shared/            ./shared/
COPY services/trip-service/ ./services/trip-service/
WORKDIR /app/services/trip-service
RUN go mod download && go build -o /trip-service .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /trip-service .
EXPOSE 50053
ENV GRPC_PORT=50053 \
    DB_HOST=postgres-trip \
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
ENTRYPOINT ["/app/trip-service"]
