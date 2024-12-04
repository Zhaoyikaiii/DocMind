package middleware

import (
	"github.com/Zhaoyikaiii/docmind/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

// AuthMiddleware handles authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Remove 'Bearer ' prefix if present
		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			utils.Logger.Warn("No authorization token provided",
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(401, gin.H{"error": "Authorization required"})
			c.Abort()
			return
		}

		// TODO: Implement actual token validation
		// For now, just pass through

		c.Next()
	}
}
