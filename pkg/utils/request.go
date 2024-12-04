package utils

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	requestCounter uint64
	startTime      = time.Now().Unix()
	hostname, _    = os.Hostname()
)

// GenerateRequestID generates a unique request ID
// 格式: hostname-timestamp-counter-uuid
func GenerateRequestID() string {
	counter := atomic.AddUint64(&requestCounter, 1)
	uuid := uuid.New().String()

	requestID := fmt.Sprintf("%s-%d-%06d-%s",
		hostname,
		startTime,
		counter,
		uuid[0:8],
	)

	return requestID
}

func ParseRequestID(requestID string) map[string]string {
	var hostname, timestamp, counter, uuid string
	fmt.Sscanf(requestID, "%s-%s-%s-%s", &hostname, &timestamp, &counter, &uuid)

	return map[string]string{
		"hostname":  hostname,
		"timestamp": timestamp,
		"counter":   counter,
		"uuid":      uuid,
	}
}

func GetRequestIDFromContext(c *gin.Context) string {
	if requestID, exists := c.Get("RequestID"); exists {
		return requestID.(string)
	}
	return "unknown"
}
