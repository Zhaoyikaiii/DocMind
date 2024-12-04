package middleware

import (
	"time"

	"github.com/Zhaoyikaiii/docmind/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger logs HTTP request details
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		requestID := utils.GenerateRequestID()
		c.Set("RequestID", requestID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		utils.Logger.Info("HTTP Request",
			zap.String("request_id", requestID),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
			zap.String("method", c.Request.Method),
			zap.Int("size", c.Writer.Size()),
		)
	}
}
