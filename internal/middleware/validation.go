package middleware

import (
	"bankapi/internal/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateRequest validates request body using struct tags
func ValidateRequest(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(model); err != nil {
			logger.Warn("Request validation failed", map[string]interface{}{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
				"ip":    c.ClientIP(),
			})

			// Format validation errors
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errors := make(map[string]string)
				for _, e := range validationErrors {
					field := strings.ToLower(e.Field())
					switch e.Tag() {
					case "required":
						errors[field] = field + " is required"
					case "email":
						errors[field] = field + " must be a valid email"
					case "min":
						errors[field] = field + " must be at least " + e.Param() + " characters"
					case "max":
						errors[field] = field + " must be at most " + e.Param() + " characters"
					case "gt":
						errors[field] = field + " must be greater than " + e.Param()
					case "gte":
						errors[field] = field + " must be greater than or equal to " + e.Param()
					case "lt":
						errors[field] = field + " must be less than " + e.Param()
					case "lte":
						errors[field] = field + " must be less than or equal to " + e.Param()
					default:
						errors[field] = field + " is invalid"
					}
				}

				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Validation failed",
					"message": "Please check your input",
					"details": errors,
					"type":    "validation_error",
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request",
					"message": err.Error(),
					"type":    "binding_error",
				})
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateQuery validates query parameters
func ValidateQuery(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindQuery(model); err != nil {
			logger.Warn("Query validation failed", map[string]interface{}{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
				"ip":    c.ClientIP(),
			})

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid query parameters",
				"message": err.Error(),
				"type":    "query_validation_error",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ValidateParams validates URL parameters
func ValidateParams(model interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindUri(model); err != nil {
			logger.Warn("Parameter validation failed", map[string]interface{}{
				"error": err.Error(),
				"path":  c.Request.URL.Path,
				"ip":    c.ClientIP(),
			})

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid URL parameters",
				"message": err.Error(),
				"type":    "parameter_validation_error",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ContentTypeValidator validates content type
func ContentTypeValidator(allowedTypes ...string) gin.HandlerFunc {
	if len(allowedTypes) == 0 {
		allowedTypes = []string{"application/json"}
	}

	return func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")
		valid := false

		for _, allowed := range allowedTypes {
			if strings.Contains(contentType, allowed) {
				valid = true
				break
			}
		}

		if !valid {
			logger.Warn("Invalid content type", map[string]interface{}{
				"content_type": contentType,
				"allowed":      allowedTypes,
				"path":         c.Request.URL.Path,
				"ip":           c.ClientIP(),
			})

			c.JSON(http.StatusUnsupportedMediaType, gin.H{
				"error":   "Unsupported content type",
				"message": "Content-Type must be one of: " + strings.Join(allowedTypes, ", "),
				"type":    "content_type_error",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
