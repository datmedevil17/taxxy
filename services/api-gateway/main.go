package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taxxy/api-gateway/config"
	"github.com/taxxy/api-gateway/routes"
)

func main() {
	cfg := config.Load()

	// ── Gin setup ────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// ── CORS – allow the React dev server (and any origin in prod) ─
	r.Use(corsMiddleware())

	// ── Mount all routes ─────────────────────────────────────
	routes.Register(r, cfg)

	// ── Listen ───────────────────────────────────────────────
	log.Printf("[GATEWAY] HTTP server listening on :%s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("[GATEWAY] run: %v", err)
	}
}

// corsMiddleware – simple CORS for development.
// In production replace with a proper middleware (e.g. gin-cors)
// and restrict origins.
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
