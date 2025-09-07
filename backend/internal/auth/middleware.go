package auth

import (
	"net/http"
	"strings"

	gin "github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT Bearer tokens and injects user identity into context.
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" || !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}
		tokenStr := strings.TrimSpace(authz[len("Bearer "):])
		claims, err := VerifyToken(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("uid", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

// WebSocketAuthMiddleware validates JWT tokens for WebSocket connections.
// Accepts token via Authorization header or query parameter.
func WebSocketAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		// First try Authorization header
		authz := c.GetHeader("Authorization")
		if authz != "" && strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			tokenStr = strings.TrimSpace(authz[len("Bearer "):])
		} else {
			// Fall back to query parameter
			tokenStr = c.Query("token")
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
			return
		}

		claims, err := VerifyToken(secret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("uid", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
