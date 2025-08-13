package audit

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	audit := router.Group("/api/v1/audit")

	// Only use auth middleware if it's provided
	if authMiddleware != nil {
		audit.Use(authMiddleware)
	}

	handler := NewHandler()

	audit.GET("/logs", handler.GetAuditLogs)
	audit.GET("/logs/id/:id", handler.GetAuditLogByID)
	audit.GET("/logs/entity/:entity_type/:entity_id", handler.GetAuditLogsByEntity)
}
