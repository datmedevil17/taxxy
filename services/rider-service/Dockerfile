FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY shared/           ./shared/
COPY services/rider-service/ ./services/rider-service/
WORKDIR /app/services/rider-service
RUN go mod download && go build -o /rider-service .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /rider-service .
EXPOSE 50052
ENV GRPC_PORT=50052 \
    DB_HOST=postgres-rider \
    TRIP_SERVICE_URL=trip-service:50053
ENTRYPOINT ["/app/rider-service"]
