FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY shared/                 ./shared/
COPY services/payment-service/ ./services/payment-service/
WORKDIR /app/services/payment-service
RUN go mod download && go build -o /payment-service .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /payment-service .
EXPOSE 50055
ENV GRPC_PORT=50055 \
    DB_HOST=postgres-payment \
    RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
ENTRYPOINT ["/app/payment-service"]
