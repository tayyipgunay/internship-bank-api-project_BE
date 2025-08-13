package balance

import (
	"bankapi/internal/db"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterRoutes(router *gin.Engine, middlewares ...gin.HandlerFunc) {
	r := router.Group("/api/v1/balances", middlewares...)
	{
		r.GET("/current", handleCurrent)
		r.GET("/historical", handleHistorical)
		r.GET("/at-time", handleAtTime)
	}
}

func handleCurrent(c *gin.Context) {
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id gerekli"})
		return
	}
	// basitçe parse et
	var userID uint
	_, err := fmt.Sscanf(userIDParam, "%d", &userID)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id geçersiz"})
		return
	}
	b, err := GetOrCreateBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bakiye getirilemedi"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func handleHistorical(c *gin.Context) {
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id gerekli"})
		return
	}
	var userID uint
	_, err := fmt.Sscanf(userIDParam, "%d", &userID)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id geçersiz"})
		return
	}
	var hist []BalanceHistory
	if err := db.DB.Where("user_id = ?", userID).Order("id DESC").Limit(200).Find(&hist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "geçmiş getirilemedi"})
		return
	}
	c.JSON(http.StatusOK, hist)
}

func handleAtTime(c *gin.Context) {
	userIDParam := c.Query("user_id")
	at := c.Query("at")
	if userIDParam == "" || at == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id ve at gerekli"})
		return
	}
	var userID uint
	_, err := fmt.Sscanf(userIDParam, "%d", &userID)
	if err != nil || userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id geçersiz"})
		return
	}
	t, err := time.Parse(time.RFC3339, at)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tarih formatı RFC3339 olmalı"})
		return
	}
	var hist BalanceHistory
	if err := db.DB.Where("user_id = ? AND created_at <= ?", userID, t).Order("created_at DESC").First(&hist).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "kayıt bulunamadı"})
		return
	}
	c.JSON(http.StatusOK, hist)
}
