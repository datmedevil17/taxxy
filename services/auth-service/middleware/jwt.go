package middleware

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ─── Claims ──────────────────────────────────────────────────

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// ─── Sign ────────────────────────────────────────────────────

// GenerateToken creates a signed JWT with the standard claims.
// secret is the value of JWT_SECRET env var.
func GenerateToken(userID, role, email, secret string, expiry time.Duration) (string, error) {
	claims := TokenClaims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "taxxy-auth",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ─── Parse ───────────────────────────────────────────────────

// ParseToken validates the signature and expiry, then returns the claims.
func ParseToken(tokenString, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
