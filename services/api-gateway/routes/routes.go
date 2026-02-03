package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/taxxy/api-gateway/config"
	"github.com/taxxy/api-gateway/handler"
)

// Register mounts every route onto the Gin engine.
func Register(r *gin.Engine, cfg config.Config) {
	gw := handler.NewGateway(cfg)

	// ── Public routes (no token required) ───────────────────
	public := r.Group("/api")
	{
		public.POST("/auth/register", gw.Register)
		public.POST("/auth/login", gw.Login)
	}

	// ── Protected routes ─────────────────────────────────────
	// Every request here goes through token validation first.
	auth := r.Group("/api")
	auth.Use(handler.AuthMiddleware(cfg.AuthServiceURL))
	{
		// Rider
		auth.POST("/rider/profile", gw.CreateRiderProfile)
		auth.GET("/rider/profile", gw.GetRiderProfile)
		auth.POST("/rider/request-ride", gw.RequestRide)
		auth.GET("/rider/trips", gw.GetRiderTrips)

		// Driver
		auth.POST("/driver/profile", gw.CreateDriverProfile)
		auth.GET("/driver/profile", gw.GetDriverProfile)
		auth.PUT("/driver/availability", gw.UpdateDriverAvailability)
		auth.POST("/driver/accept-ride", gw.AcceptRide)
		auth.POST("/driver/complete-ride", gw.CompleteRide)

		// Trip
		auth.GET("/trip/:trip_id", gw.GetTrip)

		// Payment
		auth.POST("/payment/calculate-fare", gw.CalculateFare)
		auth.POST("/payment/pay", gw.ProcessPayment)
	}

	// ── Health check (used by k8s liveness/readiness probes) ─
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}
