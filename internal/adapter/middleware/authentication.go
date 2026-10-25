package middleware

import (
	"demo/internal/core/port"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate(tokenManager port.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Authorization: Bearer <TOKEN>
		tokenString, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		userID, role, err := tokenManager.ValidateToken(tokenString)
		if err != nil || role == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token claims"})
			return
		}

		c.Set("userID", userID)
		c.Set("role", role)
		c.Next()
	}
}
