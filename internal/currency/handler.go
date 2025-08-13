package currency

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *CurrencyService
}

func NewHandler() *Handler {
	return &Handler{
		service: NewCurrencyService(),
	}
}

// ConvertCurrencyRequest represents the request to convert currency
type ConvertCurrencyRequest struct {
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	FromCurrency string  `json:"from_currency" binding:"required"`
	ToCurrency   string  `json:"to_currency" binding:"required"`
}

// ConvertCurrency converts an amount from one currency to another
func (h *Handler) ConvertCurrency(c *gin.Context) {
	var req ConvertCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update rates from external source
	if err := h.service.UpdateRates(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exchange rates"})
		return
	}

	convertedAmount, err := h.service.Convert(req.Amount, req.FromCurrency, req.ToCurrency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original_amount":  req.Amount,
		"from_currency":    req.FromCurrency,
		"to_currency":      req.ToCurrency,
		"converted_amount": convertedAmount,
	})
}

// GetExchangeRate returns the exchange rate between two currencies
func (h *Handler) GetExchangeRate(c *gin.Context) {
	fromCurrency := c.Param("from")
	toCurrency := c.Param("to")

	// Update rates from external source
	if err := h.service.UpdateRates(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exchange rates"})
		return
	}

	rate, err := h.service.GetExchangeRate(fromCurrency, toCurrency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"from_currency": fromCurrency,
		"to_currency":   toCurrency,
		"rate":          rate,
	})
}

// GetSupportedCurrencies returns list of supported currencies
func (h *Handler) GetSupportedCurrencies(c *gin.Context) {
	currencies := h.service.GetSupportedCurrencies()
	c.JSON(http.StatusOK, gin.H{
		"currencies": currencies,
	})
}

// UpdateExchangeRate updates the exchange rate between two currencies
func (h *Handler) UpdateExchangeRate(c *gin.Context) {
	fromCurrency := c.Param("from")
	toCurrency := c.Param("to")

	var req struct {
		Rate float64 `json:"rate" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.service.SetExchangeRate(fromCurrency, toCurrency, req.Rate)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Exchange rate updated successfully",
		"from_currency": fromCurrency,
		"to_currency":   toCurrency,
		"rate":          req.Rate,
	})
}
