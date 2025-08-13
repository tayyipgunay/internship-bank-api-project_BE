package transaction

import (
	"bankapi/internal/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(router *gin.Engine, middlewares ...gin.HandlerFunc) {
	r := router.Group("/api/v1/transactions", middlewares...)
	{
		r.POST("/credit", handleCredit)
		r.POST("/debit", handleDebit)
		r.POST("/transfer", handleTransfer)
		r.GET("/history", handleHistory)
		r.GET("/:id", handleGetByID)
	}
}

func handleCredit(c *gin.Context) {
	var req CreateCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}
	tx, err := ApplyCredit(req.UserID, req.AmountCents)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "transaction": tx})
		return
	}
	c.JSON(http.StatusCreated, tx)
}

func handleDebit(c *gin.Context) {
	var req CreateDebitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}
	tx, err := ApplyDebit(req.UserID, req.AmountCents)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "transaction": tx})
		return
	}
	c.JSON(http.StatusCreated, tx)
}

func handleTransfer(c *gin.Context) {
	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri"})
		return
	}
	tx, err := ApplyTransfer(req.FromUserID, req.ToUserID, req.AmountCents)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "transaction": tx})
		return
	}
	c.JSON(http.StatusCreated, tx)
}

func handleHistory(c *gin.Context) {
	// basic: list latest transactions; can add filters later
	var txs []Transaction
	if err := db.DB.Order("id DESC").Limit(100).Find(&txs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "geçmiş getirilemedi"})
		return
	}
	c.JSON(http.StatusOK, txs)
}

func handleGetByID(c *gin.Context) {
	var tx Transaction
	id := c.Param("id")
	if err := db.DB.First(&tx, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "işlem bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, tx)
}
