package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/taxxy/shared/proto/auth"
)

// ─── Context keys ────────────────────────────────────────────

const (
	CtxUserID = "user_id"
	CtxRole   = "role"
	CtxEmail  = "email"
)

// ─── AuthMiddleware ──────────────────────────────────────────
// Extracts "Authorization: Bearer <token>", validates via auth-service,
// then sets user_id / role / email on the Gin context.

func AuthMiddleware(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// ── Dial auth-service ──────────────────────────────
		conn, err := grpc.NewClient(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth service unavailable"})
			return
		}
		defer conn.Close()

		client := authpb.NewAuthServiceClient(conn)
		resp, err := client.ValidateToken(c.Request.Context(), &authpb.ValidateTokenRequest{Token: token})
		if err != nil || !resp.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// ── Inject into context ────────────────────────────
		c.Set(CtxUserID, resp.UserId)
		c.Set(CtxRole, resp.Role)
		c.Set(CtxEmail, resp.Email)

		c.Next()
	}
}

// ─── Helper to pull user_id from context ─────────────────────

func GetUserID(c *gin.Context) string {
	v, _ := c.Get(CtxUserID)
	return v.(string)
}
