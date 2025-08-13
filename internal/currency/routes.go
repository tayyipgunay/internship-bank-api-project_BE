package currency

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	curr := router.Group("/api/v1/currency")

	// Only use auth middleware if it's provided
	if authMiddleware != nil {
		curr.Use(authMiddleware)
	}

	handler := NewHandler()

	curr.POST("/convert", handler.ConvertCurrency)
	curr.GET("/rates/:from/:to", handler.GetExchangeRate)
	curr.GET("/currencies", handler.GetSupportedCurrencies)
	curr.PUT("/rates/:from/:to", handler.UpdateExchangeRate)
}
