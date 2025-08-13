package currency

import (
	"fmt"
	"sync"
	"time"
)

// Currency represents a currency with its code and exchange rates
type Currency struct {
	Code         string    `json:"code" gorm:"primaryKey;size:3"`
	Name         string    `json:"name" gorm:"size:100;not null"`
	Symbol       string    `json:"symbol" gorm:"size:5"`
	ExchangeRate float64   `json:"exchange_rate" gorm:"not null"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ExchangeRate represents an exchange rate between two currencies
type ExchangeRate struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	FromCurrency string    `json:"from_currency" gorm:"size:3;index;not null"`
	ToCurrency   string    `json:"to_currency" gorm:"size:3;index;not null"`
	Rate         float64   `json:"rate" gorm:"not null"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CurrencyService manages currency operations
type CurrencyService struct {
	rates map[string]float64
	mutex sync.RWMutex
}

// NewCurrencyService creates a new currency service
func NewCurrencyService() *CurrencyService {
	println("💱 Currency service oluşturuluyor...")

	service := &CurrencyService{
		rates: make(map[string]float64),
	}

	println("✅ Currency service oluşturuldu")
	return service
}

// Convert converts an amount from one currency to another
func (cs *CurrencyService) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	rateKey := fmt.Sprintf("%s_%s", fromCurrency, toCurrency)
	rate, exists := cs.rates[rateKey]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", fromCurrency, toCurrency)
	}

	return amount * rate, nil
}

// SetExchangeRate sets an exchange rate between two currencies
func (cs *CurrencyService) SetExchangeRate(fromCurrency, toCurrency string, rate float64) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	rateKey := fmt.Sprintf("%s_%s", fromCurrency, toCurrency)
	cs.rates[rateKey] = rate

	// Set reverse rate
	reverseKey := fmt.Sprintf("%s_%s", toCurrency, fromCurrency)
	cs.rates[reverseKey] = 1.0 / rate
}

// GetExchangeRate gets the exchange rate between two currencies
func (cs *CurrencyService) GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	rateKey := fmt.Sprintf("%s_%s", fromCurrency, toCurrency)
	rate, exists := cs.rates[rateKey]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", fromCurrency, toCurrency)
	}

	return rate, nil
}

// UpdateRates updates exchange rates from external source
func (cs *CurrencyService) UpdateRates() error {
	// In a real implementation, you would fetch rates from an external API
	// For now, we'll set some sample rates

	cs.SetExchangeRate("USD", "EUR", 0.85)
	cs.SetExchangeRate("USD", "TRY", 30.0)
	cs.SetExchangeRate("EUR", "TRY", 35.0)

	return nil
}

// GetSupportedCurrencies returns list of supported currencies
func (cs *CurrencyService) GetSupportedCurrencies() []string {
	return []string{"USD", "EUR", "TRY"}
}
