package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/Zhaoyikaiii/docmind/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery handles panic recovery
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()

				utils.Logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.ByteString("stack", stack),
					zap.String("path", c.Request.URL.Path),
					zap.String("request_id", c.GetString("RequestID")),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal Server Error",
					"request_id": c.GetString("RequestID"),
				})
			}
		}()
		c.Next()
	}
}
