package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware records metrics for HTTP requests
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record start time
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get route pattern
		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}

		// Record metrics
		RecordAPICall(
			c.Request.Context(),
			c.Request.Method,
			route,
			c.Writer.Status(),
			duration,
		)
	}
}
