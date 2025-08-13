package worker

import (
	"sync/atomic"
	"time"
)

// TransactionStats tracks transaction statistics using atomic operations
type TransactionStats struct {
	TotalTransactions      int64
	SuccessfulTransactions int64
	FailedTransactions     int64
	PendingTransactions    int64
	TotalAmountCents       int64
	LastTransactionTime    int64 // Unix timestamp
}

// GetStats returns a copy of current statistics
func (ts *TransactionStats) GetStats() TransactionStats {
	return TransactionStats{
		TotalTransactions:      atomic.LoadInt64(&ts.TotalTransactions),
		SuccessfulTransactions: atomic.LoadInt64(&ts.SuccessfulTransactions),
		FailedTransactions:     atomic.LoadInt64(&ts.FailedTransactions),
		PendingTransactions:    atomic.LoadInt64(&ts.PendingTransactions),
		TotalAmountCents:       atomic.LoadInt64(&ts.TotalAmountCents),
		LastTransactionTime:    atomic.LoadInt64(&ts.LastTransactionTime),
	}
}

// IncrementTotal increments total transaction count
func (ts *TransactionStats) IncrementTotal() {
	atomic.AddInt64(&ts.TotalTransactions, 1)
}

// IncrementSuccessful increments successful transaction count
func (ts *TransactionStats) IncrementSuccessful() {
	atomic.AddInt64(&ts.SuccessfulTransactions, 1)
}

// IncrementFailed increments failed transaction count
func (ts *TransactionStats) IncrementFailed() {
	atomic.AddInt64(&ts.FailedTransactions, 1)
}

// IncrementPending increments pending transaction count
func (ts *TransactionStats) IncrementPending() {
	atomic.AddInt64(&ts.PendingTransactions, 1)
}

// DecrementPending decrements pending transaction count
func (ts *TransactionStats) DecrementPending() {
	atomic.AddInt64(&ts.PendingTransactions, -1)
}

// AddAmount adds amount to total
func (ts *TransactionStats) AddAmount(amountCents int64) {
	atomic.AddInt64(&ts.TotalAmountCents, amountCents)
}

// UpdateLastTransactionTime updates the last transaction timestamp
func (ts *TransactionStats) UpdateLastTransactionTime() {
	atomic.StoreInt64(&ts.LastTransactionTime, time.Now().Unix())
}

// GetSuccessRate returns the success rate as a percentage
func (ts *TransactionStats) GetSuccessRate() float64 {
	total := atomic.LoadInt64(&ts.TotalTransactions)
	if total == 0 {
		return 0.0
	}
	successful := atomic.LoadInt64(&ts.SuccessfulTransactions)
	return float64(successful) / float64(total) * 100.0
}

// GetAverageAmount returns the average transaction amount
func (ts *TransactionStats) GetAverageAmount() float64 {
	total := atomic.LoadInt64(&ts.TotalTransactions)
	if total == 0 {
		return 0.0
	}
	amount := atomic.LoadInt64(&ts.TotalAmountCents)
	return float64(amount) / float64(total) / 100.0 // Convert cents to TL
}

// Reset resets all statistics to zero
func (ts *TransactionStats) Reset() {
	atomic.StoreInt64(&ts.TotalTransactions, 0)
	atomic.StoreInt64(&ts.SuccessfulTransactions, 0)
	atomic.StoreInt64(&ts.FailedTransactions, 0)
	atomic.StoreInt64(&ts.PendingTransactions, 0)
	atomic.StoreInt64(&ts.TotalAmountCents, 0)
	atomic.StoreInt64(&ts.LastTransactionTime, 0)
}
