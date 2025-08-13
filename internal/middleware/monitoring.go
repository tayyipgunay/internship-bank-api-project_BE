package middleware

import (
	"bankapi/internal/logger"
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// PerformanceMonitor tracks request performance metrics
func PerformanceMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		status := c.Writer.Status()

		// Log performance metrics
		logger.Info("Request Performance", map[string]interface{}{
			"method":         c.Request.Method,
			"path":           c.Request.URL.Path,
			"status":         status,
			"duration_ms":    duration.Milliseconds(),
			"duration_ns":    duration.Nanoseconds(),
			"content_length": c.Writer.Size(),
			"ip":             c.ClientIP(),
		})

		// Track slow requests (>1 second)
		if duration > time.Second {
			logger.Warn("Slow Request Detected", map[string]interface{}{
				"method":      c.Request.Method,
				"path":        c.Request.URL.Path,
				"duration_ms": duration.Milliseconds(),
				"ip":          c.ClientIP(),
			})
		}

		// Track error responses
		if status >= 400 {
			logger.Error("Error Response", fmt.Errorf("HTTP %d", status), map[string]interface{}{
				"method":      c.Request.Method,
				"path":        c.Request.URL.Path,
				"status":      status,
				"duration_ms": duration.Milliseconds(),
				"ip":          c.ClientIP(),
			})
		}
	}
}

// MemoryUsage tracks memory usage for requests
func MemoryUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		// Log memory usage before request
		logger.Debug("Memory Before Request", map[string]interface{}{
			"alloc_mb":       bToMb(m.Alloc),
			"total_alloc_mb": bToMb(m.TotalAlloc),
			"sys_mb":         bToMb(m.Sys),
			"num_gc":         m.NumGC,
		})

		c.Next()

		// Log memory usage after request
		runtime.ReadMemStats(&m)
		logger.Debug("Memory After Request", map[string]interface{}{
			"alloc_mb":       bToMb(m.Alloc),
			"total_alloc_mb": bToMb(m.TotalAlloc),
			"sys_mb":         bToMb(m.Sys),
			"num_gc":         m.NumGC,
		})
	}
}

// bToMb converts bytes to megabytes
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
