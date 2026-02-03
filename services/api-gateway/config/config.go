package config

import "os"

type Config struct {
	HTTPPort           string
	AuthServiceURL     string
	RiderServiceURL    string
	DriverServiceURL   string
	TripServiceURL     string
	PaymentServiceURL  string
}

func Load() Config {
	return Config{
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "auth-service:50051"),
		RiderServiceURL:   getEnv("RIDER_SERVICE_URL", "rider-service:50052"),
		DriverServiceURL:  getEnv("DRIVER_SERVICE_URL", "driver-service:50054"),
		TripServiceURL:    getEnv("TRIP_SERVICE_URL", "trip-service:50053"),
		PaymentServiceURL: getEnv("PAYMENT_SERVICE_URL", "payment-service:50055"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
