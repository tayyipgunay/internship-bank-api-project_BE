package middleware

import (
	"bankapi/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logger.Error("Panic recovered", nil, map[string]interface{}{
				"error": err,
				"path":  c.Request.URL.Path,
				"ip":    c.ClientIP(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Something went wrong",
			})
		} else {
			logger.Error("Unknown panic recovered", nil, map[string]interface{}{
				"recovered": recovered,
				"path":      c.Request.URL.Path,
				"ip":        c.ClientIP(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Unknown error occurred",
			})
		}
	})
}

// ValidationError formats validation errors
func ValidationError(c *gin.Context, err error) {
	logger.Warn("Validation error", map[string]interface{}{
		"error": err.Error(),
		"path":  c.Request.URL.Path,
		"ip":    c.ClientIP(),
	})

	c.JSON(http.StatusBadRequest, gin.H{
		"error":   "Validation failed",
		"message": err.Error(),
		"type":    "validation_error",
	})
}

// NotFoundError formats not found errors
func NotFoundError(c *gin.Context, resource string) {
	logger.Warn("Resource not found", map[string]interface{}{
		"resource": resource,
		"path":     c.Request.URL.Path,
		"ip":       c.ClientIP(),
	})

	c.JSON(http.StatusNotFound, gin.H{
		"error":   "Resource not found",
		"message": resource + " not found",
		"type":    "not_found",
	})
}

// UnauthorizedError formats unauthorized errors
func UnauthorizedError(c *gin.Context, message string) {
	logger.Warn("Unauthorized access", map[string]interface{}{
		"message": message,
		"path":    c.Request.URL.Path,
		"ip":      c.ClientIP(),
	})

	c.JSON(http.StatusUnauthorized, gin.H{
		"error":   "Unauthorized",
		"message": message,
		"type":    "unauthorized",
	})
}

// ForbiddenError formats forbidden errors
func ForbiddenError(c *gin.Context, message string) {
	logger.Warn("Forbidden access", map[string]interface{}{
		"message": message,
		"path":    c.Request.URL.Path,
		"ip":      c.ClientIP(),
	})

	c.JSON(http.StatusForbidden, gin.H{
		"error":   "Forbidden",
		"message": message,
		"type":    "forbidden",
	})
}
