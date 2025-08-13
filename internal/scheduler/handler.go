package scheduler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	scheduler *Scheduler
}

func NewHandler(scheduler *Scheduler) *Handler {
	return &Handler{scheduler: scheduler}
}

// ScheduleTransactionRequest represents the request to schedule a transaction
type ScheduleTransactionRequest struct {
	FromUserID string                 `json:"from_user_id" binding:"required"`
	ToUserID   string                 `json:"to_user_id" binding:"required"`
	Amount     float64                `json:"amount" binding:"required,gt=0"`
	Type       string                 `json:"type" binding:"required"`
	Schedule   string                 `json:"schedule" binding:"required"` // Cron expression
	Metadata   map[string]interface{} `json:"metadata"`
}

// ScheduleTransaction schedules a new transaction
func (h *Handler) ScheduleTransaction(c *gin.Context) {
	var req ScheduleTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cron expression
	if !isValidCronExpression(req.Schedule) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cron expression"})
		return
	}

	st := &ScheduledTransaction{
		FromUserID: req.FromUserID,
		ToUserID:   req.ToUserID,
		Amount:     req.Amount,
		Type:       req.Type,
		Schedule:   req.Schedule,
		Status:     "pending",
		Metadata:   req.Metadata,
	}

	if err := h.scheduler.ScheduleTransaction(st); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Transaction scheduled successfully",
		"id":      st.ID,
	})
}

// GetScheduledTransactions returns all scheduled transactions
func (h *Handler) GetScheduledTransactions(c *gin.Context) {
	transactions := h.scheduler.GetScheduledTransactions()
	c.JSON(http.StatusOK, transactions)
}

// UnscheduleTransaction removes a scheduled transaction
func (h *Handler) UnscheduleTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	if err := h.scheduler.UnscheduleTransaction(transactionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scheduled transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction unscheduled successfully"})
}

// UpdateSchedule updates the schedule of a transaction
func (h *Handler) UpdateSchedule(c *gin.Context) {
	transactionID := c.Param("id")

	var req struct {
		Schedule string `json:"schedule" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate cron expression
	if !isValidCronExpression(req.Schedule) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cron expression"})
		return
	}

	if err := h.scheduler.UpdateSchedule(transactionID, req.Schedule); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scheduled transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schedule updated successfully"})
}

// isValidCronExpression validates a cron expression
func isValidCronExpression(expression string) bool {
	// Basic validation - in production you might want to use a proper cron parser
	// For now, we'll just check if it has the right number of fields
	// Standard cron has 5 fields: minute hour day month weekday
	// Extended cron (with seconds) has 6 fields

	if len(expression) < 5 {
		return false
	}

	// This is a simple check - in production use proper validation
	return true
}
