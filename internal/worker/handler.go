package worker

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(router *gin.Engine, p *Processor) {
	r := router.Group("/api/v1/queue")
	{
		r.POST("/credit", func(c *gin.Context) {
			var req struct {
				UserID uint  `json:"user_id"`
				Amount int64 `json:"amount_cents"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
				return
			}
			var j Job
			j.Kind = "credit"
			j.Credit.UserID = req.UserID
			j.Credit.Amount = req.Amount
			p.Enqueue(j)
			c.Status(http.StatusAccepted)
		})
		r.POST("/debit", func(c *gin.Context) {
			var req struct {
				UserID uint  `json:"user_id"`
				Amount int64 `json:"amount_cents"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
				return
			}
			var j Job
			j.Kind = "debit"
			j.Debit.UserID = req.UserID
			j.Debit.Amount = req.Amount
			p.Enqueue(j)
			c.Status(http.StatusAccepted)
		})
		r.POST("/transfer", func(c *gin.Context) {
			var req struct {
				FromID uint  `json:"from_user_id"`
				ToID   uint  `json:"to_user_id"`
				Amount int64 `json:"amount_cents"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
				return
			}
			var j Job
			j.Kind = "transfer"
			j.Transfer.FromID = req.FromID
			j.Transfer.ToID = req.ToID
			j.Transfer.Amount = req.Amount
			p.Enqueue(j)
			c.Status(http.StatusAccepted)
		})
		r.GET("/stats", func(c *gin.Context) {
			ok, ng := p.Stats()
			c.JSON(http.StatusOK, gin.H{"ok": ok, "failed": ng})
		})
	}
}
