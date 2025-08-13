package scheduler

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc, scheduler *Scheduler) {
	sched := router.Group("/api/v1/scheduler")

	// Only use auth middleware if it's provided
	if authMiddleware != nil {
		sched.Use(authMiddleware)
	}

	handler := NewHandler(scheduler)

	sched.POST("/transactions", handler.ScheduleTransaction)
	sched.GET("/transactions", handler.GetScheduledTransactions)
	sched.DELETE("/transactions/:id", handler.UnscheduleTransaction)
	sched.PUT("/transactions/:id/schedule", handler.UpdateSchedule)
}
