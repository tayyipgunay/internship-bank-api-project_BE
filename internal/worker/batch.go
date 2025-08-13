package worker

import (
	"bankapi/internal/logger"
	"bankapi/internal/transaction"
	"context"
	"fmt"
	"sync"
	"time"
)

// BatchProcessor handles batch transaction processing
type BatchProcessor struct {
	workerPool   *WorkerPool
	batchSize    int
	batchTimeout time.Duration
	stats        *TransactionStats
	mu           sync.RWMutex
}

// BatchTransaction represents a batch of transactions
type BatchTransaction struct {
	ID           string
	Transactions []transaction.Transaction
	Status       BatchStatus
	CreatedAt    time.Time
	CompletedAt  *time.Time
	Error        string
}

type BatchStatus string

const (
	BatchStatusPending    BatchStatus = "pending"
	BatchStatusProcessing BatchStatus = "processing"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusFailed     BatchStatus = "failed"
)

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(workerPool *WorkerPool, batchSize int, batchTimeout time.Duration) *BatchProcessor {
	return &BatchProcessor{
		workerPool:   workerPool,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
		stats:        &TransactionStats{},
	}
}

// ProcessBatch processes a batch of transactions
func (bp *BatchProcessor) ProcessBatch(ctx context.Context, transactions []transaction.Transaction) (*BatchTransaction, error) {
	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions to process")
	}

	if len(transactions) > bp.batchSize {
		return nil, fmt.Errorf("batch size %d exceeds maximum %d", len(transactions), bp.batchSize)
	}

	batch := &BatchTransaction{
		ID:           generateBatchID(),
		Transactions: transactions,
		Status:       BatchStatusPending,
		CreatedAt:    time.Now(),
	}

	// Start processing in background
	go bp.processBatchAsync(ctx, batch)

	return batch, nil
}

// processBatchAsync processes the batch asynchronously
func (bp *BatchProcessor) processBatchAsync(ctx context.Context, batch *BatchTransaction) {
	bp.mu.Lock()
	batch.Status = BatchStatusProcessing
	bp.mu.Unlock()

	// Create context with timeout
	batchCtx, cancel := context.WithTimeout(ctx, bp.batchTimeout)
	defer cancel()

	// Process transactions concurrently
	var wg sync.WaitGroup
	results := make(chan *BatchResult, len(batch.Transactions))
	errors := make(chan error, len(batch.Transactions))

	// Submit transactions to worker pool
	for _, txn := range batch.Transactions {
		wg.Add(1)
		go func(t transaction.Transaction) {
			defer wg.Done()

			// Submit to worker pool
			result, err := bp.workerPool.SubmitTransaction(batchCtx, t)
			if err != nil {
				errors <- err
				return
			}

			results <- result
		}(txn)
	}

	// Wait for all transactions to complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	var successfulCount int
	var failedCount int
	var totalAmount int64

	// Process results
	for result := range results {
		if result != nil && result.Success {
			successfulCount++
			totalAmount += result.Transaction.AmountCents
		} else {
			failedCount++
		}
	}

	// Process errors
	for err := range errors {
		if err != nil {
			failedCount++
			logger.Error("Batch transaction failed", err, map[string]interface{}{
				"batch_id": batch.ID,
			})
		}
	}

	// Update batch status
	bp.mu.Lock()
	if failedCount == 0 {
		batch.Status = BatchStatusCompleted
	} else if successfulCount == 0 {
		batch.Status = BatchStatusFailed
		batch.Error = "All transactions failed"
	} else {
		batch.Status = BatchStatusCompleted
		batch.Error = fmt.Sprintf("%d transactions failed", failedCount)
	}
	batch.CompletedAt = &time.Time{}
	*batch.CompletedAt = time.Now()
	bp.mu.Unlock()

	// Update statistics
	bp.stats.IncrementTotal()
	if batch.Status == BatchStatusCompleted {
		bp.stats.IncrementSuccessful()
	} else {
		bp.stats.IncrementFailed()
	}
	bp.stats.AddAmount(totalAmount)
	bp.stats.UpdateLastTransactionTime()

	logger.Info("Batch processing completed", map[string]interface{}{
		"batch_id":        batch.ID,
		"total":           len(batch.Transactions),
		"successful":      successfulCount,
		"failed":          failedCount,
		"status":          batch.Status,
		"processing_time": time.Since(batch.CreatedAt),
	})
}

// GetBatchStatus returns the current status of a batch
func (bp *BatchProcessor) GetBatchStatus(batchID string) (*BatchTransaction, error) {
	// This is a simplified implementation
	// In a real system, you'd store batches in a database
	return nil, fmt.Errorf("batch status tracking not implemented")
}

// GetStats returns batch processing statistics
func (bp *BatchProcessor) GetStats() *TransactionStats {
	return bp.stats
}

// generateBatchID generates a unique batch ID
func generateBatchID() string {
	return fmt.Sprintf("batch_%d_%s", time.Now().UnixNano(), randomString(8))
}

// randomString generates a random string of given length
func randomString(length int) string {
	println("🎲 Random string oluşturuluyor, uzunluk:", length)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	result := string(b)
	println("🎲 Random string oluşturuldu:", result)
	return result
}

// BatchResult represents the result of processing a single transaction
type BatchResult struct {
	Transaction transaction.Transaction
	Success     bool
	Error       string
	ProcessedAt time.Time
}
