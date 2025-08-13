package worker

import (
	"bankapi/internal/transaction"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	Kind   string
	Credit struct {
		UserID uint
		Amount int64
	}
	Debit struct {
		UserID uint
		Amount int64
	}
	Transfer struct {
		FromID uint
		ToID   uint
		Amount int64
	}
}

type Processor struct {
	queue       chan Job
	numWorkers  int
	processedOK atomic.Int64
	processedNG atomic.Int64
}

// WorkerPool handles concurrent transaction processing
type WorkerPool struct {
	workers  int
	jobQueue chan *TransactionJob
	stats    *TransactionStats
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// TransactionJob represents a transaction job
type TransactionJob struct {
	ID          string
	Transaction transaction.Transaction
	ResultChan  chan *TransactionResult
}

// TransactionResult represents the result of processing a transaction
type TransactionResult struct {
	ID          string
	Success     bool
	Error       string
	ProcessedAt time.Time
}

func NewProcessor(buffer, workers int) *Processor {
	return &Processor{queue: make(chan Job, buffer), numWorkers: workers}
}

func (p *Processor) Enqueue(j Job) { p.queue <- j }

func (p *Processor) ProcessForever() {
	println("🚀 Worker pool başlatılıyor, worker sayısı:", p.numWorkers)

	for i := 0; i < p.numWorkers; i++ {
		workerID := i
		go func() {
			println("👷 Worker", workerID, "başlatıldı")
			for j := range p.queue {
				println("🔧 Worker", workerID, "işi işliyor, tip:", j.Kind)

				var err error
				switch j.Kind {
				case "credit":
					println("💳 Kredi işlemi işleniyor, kullanıcı:", j.Credit.UserID, "miktar:", j.Credit.Amount)
					_, err = transaction.ApplyCredit(j.Credit.UserID, j.Credit.Amount)
				case "debit":
					println("💸 Debit işlemi işleniyor, kullanıcı:", j.Debit.UserID, "miktar:", j.Debit.Amount)
					_, err = transaction.ApplyDebit(j.Debit.UserID, j.Debit.Amount)
				case "transfer":
					println("🔄 Transfer işlemi işleniyor, from:", j.Transfer.FromID, "to:", j.Transfer.ToID, "miktar:", j.Transfer.Amount)
					_, err = transaction.ApplyTransfer(j.Transfer.FromID, j.Transfer.ToID, j.Transfer.Amount)
				default:
					println("❌ Bilinmeyen iş tipi:", j.Kind)
					err = fmt.Errorf("unknown job kind: %s", j.Kind)
				}

				if err != nil {
					println("❌ Worker", workerID, "işi başarısız:", err.Error())
					p.processedNG.Add(1)
				} else {
					println("✅ Worker", workerID, "işi başarılı")
					p.processedOK.Add(1)
				}
			}
			println("👷 Worker", workerID, "durduruldu")
		}()
	}
	println("✅ Worker pool başlatıldı")
}

func (p *Processor) Stats() (ok, ng int64) {
	return p.processedOK.Load(), p.processedNG.Load()
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers:  workers,
		jobQueue: make(chan *TransactionJob, workers*10), // Buffer size
		stats:    &TransactionStats{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	println("🚀 Worker pool başlatılıyor, worker sayısı:", wp.workers)

	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		workerID := i
		go func() {
			println("👷 Worker", workerID, "başlatıldı")
			wp.worker(workerID)
		}()
	}
	println("✅ Worker pool başlatıldı")
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	println("🛑 Worker pool durduruluyor...")
	wp.cancel()
	close(wp.jobQueue)
	wp.wg.Wait()
	println("✅ Worker pool durduruldu")
}

// SubmitTransaction submits a transaction for processing
func (wp *WorkerPool) SubmitTransaction(ctx context.Context, txn transaction.Transaction) (*TransactionResult, error) {
	println("📝 Transaction worker pool'a gönderiliyor, ID:", txn.ID)

	job := &TransactionJob{
		ID:          generateJobID(),
		Transaction: txn,
		ResultChan:  make(chan *TransactionResult, 1),
	}

	select {
	case wp.jobQueue <- job:
		println("✅ Transaction başarıyla gönderildi, job ID:", job.ID)
	case <-ctx.Done():
		println("❌ Context iptal edildi")
		return nil, ctx.Err()
	default:
		println("❌ Worker pool kuyruğu dolu")
		return nil, fmt.Errorf("worker pool queue is full")
	}

	// Wait for result
	println("⏳ Transaction sonucu bekleniyor...")
	select {
	case result := <-job.ResultChan:
		println("✅ Transaction sonucu alındı, başarılı:", result.Success)
		return result, nil
	case <-ctx.Done():
		println("❌ Context timeout, sonuç alınamadı")
		return nil, ctx.Err()
	}
}

// worker processes jobs from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok {
				println("👷 Worker", id, "kuyruk kapandı, durduruluyor")
				return // Channel closed
			}

			println("🔧 Worker", id, "transaction işliyor, job ID:", job.ID)

			// Process the transaction
			result := wp.processTransaction(job)

			// Send result
			select {
			case job.ResultChan <- result:
				println("✅ Worker", id, "sonucu gönderdi")
			default:
				println("⚠️ Worker", id, "sonuç kanalı dolu")
			}

		case <-wp.ctx.Done():
			println("👷 Worker", id, "context iptal edildi, durduruluyor")
			return
		}
	}
}

// processTransaction processes a single transaction
func (wp *WorkerPool) processTransaction(job *TransactionJob) *TransactionResult {
	println("🔧 Transaction işleniyor, job ID:", job.ID)

	wp.stats.IncrementTotal()
	wp.stats.IncrementPending()

	// Simulate transaction processing
	// In a real system, this would interact with the database
	success := true
	var errMsg string

	// Update statistics
	if success {
		wp.stats.IncrementSuccessful()
		wp.stats.AddAmount(job.Transaction.AmountCents)
		println("✅ Transaction başarılı, miktar:", job.Transaction.AmountCents)
	} else {
		wp.stats.IncrementFailed()
		errMsg = "Transaction processing failed"
		println("❌ Transaction başarısız")
	}

	wp.stats.DecrementPending()
	wp.stats.UpdateLastTransactionTime()

	result := &TransactionResult{
		ID:          job.ID,
		Success:     success,
		Error:       errMsg,
		ProcessedAt: time.Now(),
	}

	println("✅ Transaction işlemi tamamlandı, sonuç:", result.Success)
	return result
}

// generateJobID generates a unique job ID
func generateJobID() string {
	println("🆔 Job ID oluşturuluyor...")
	id := fmt.Sprintf("job_%d_%s", time.Now().UnixNano(), randomString(8))
	println("🆔 Job ID oluşturuldu:", id)
	return id
}
