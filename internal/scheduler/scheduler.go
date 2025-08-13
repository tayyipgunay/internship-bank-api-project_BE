package scheduler

import (
	"bankapi/internal/events"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// ScheduledTransaction represents a transaction that will be executed at a specific time
type ScheduledTransaction struct {
	ID         string                 `json:"id"`
	FromUserID string                 `json:"from_user_id"`
	ToUserID   string                 `json:"to_user_id"`
	Amount     float64                `json:"amount"`
	Type       string                 `json:"type"`
	Schedule   string                 `json:"schedule"` // Cron expression
	Status     string                 `json:"status"`   // pending, completed, failed
	CreatedAt  time.Time              `json:"created_at"`
	ExecuteAt  time.Time              `json:"execute_at"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// Scheduler manages scheduled transactions
type Scheduler struct {
	cron     *cron.Cron
	entries  map[string]cron.EntryID
	mutex    sync.RWMutex
	eventBus events.EventBus
}

// NewScheduler creates a new scheduler instance
func NewScheduler(eventBus events.EventBus) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		entries:  make(map[string]cron.EntryID),
		eventBus: eventBus,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	println("⏰ Scheduler başlatılıyor...")
	s.cron.Start()
	println("✅ Scheduler başlatıldı")
	log.Println("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	println("⏰ Scheduler durduruluyor...")
	s.cron.Stop()
	println("✅ Scheduler durduruldu")
	log.Println("Scheduler stopped")
}

// ScheduleTransaction schedules a new transaction
func (s *Scheduler) ScheduleTransaction(st *ScheduledTransaction) error {
	println("⏰ Transaction planlanıyor, ID:", st.ID)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Validate required fields
	if st.ID == "" {
		println("❌ Transaction ID boş")
		return fmt.Errorf("transaction ID is required")
	}

	if st.Schedule == "" {
		println("❌ Schedule boş")
		return fmt.Errorf("schedule is required")
	}

	if st.Amount <= 0 {
		println("❌ Geçersiz miktar:", st.Amount)
		return fmt.Errorf("amount must be positive")
	}

	// Parse cron expression
	entryID, err := s.cron.AddFunc(st.Schedule, func() {
		println("🚀 Planlanan transaction çalıştırılıyor, ID:", st.ID)
		s.executeScheduledTransaction(st)
	})
	if err != nil {
		println("❌ Transaction planlanamadı:", err.Error())
		return fmt.Errorf("failed to schedule transaction: %w", err)
	}

	s.entries[st.ID] = entryID
	println("✅ Transaction başarıyla planlandı, ID:", st.ID, "Entry ID:", entryID)
	log.Printf("Scheduled transaction %s for execution", st.ID)

	return nil
}

// UnscheduleTransaction removes a scheduled transaction
func (s *Scheduler) UnscheduleTransaction(transactionID string) error {
	println("⏰ Transaction planı kaldırılıyor, ID:", transactionID)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	entryID, exists := s.entries[transactionID]
	if !exists {
		println("❌ Planlanan transaction bulunamadı:", transactionID)
		return fmt.Errorf("scheduled transaction not found: %s", transactionID)
	}

	s.cron.Remove(entryID)
	delete(s.entries, transactionID)
	println("✅ Transaction planı kaldırıldı, ID:", transactionID)
	log.Printf("Unscheduled transaction %s", transactionID)

	return nil
}

// executeScheduledTransaction executes a scheduled transaction
func (s *Scheduler) executeScheduledTransaction(st *ScheduledTransaction) {
	println("🚀 Planlanan transaction çalıştırılıyor, ID:", st.ID)

	// Create transaction event
	event := events.NewEvent("transaction.scheduled", st.ID, map[string]interface{}{
		"from_user_id": st.FromUserID,
		"to_user_id":   st.ToUserID,
		"amount":       st.Amount,
		"type":         st.Type,
		"metadata":     st.Metadata,
	})

	// Publish event
	if err := s.eventBus.Publish(event); err != nil {
		println("❌ Planlanan transaction event yayınlanamadı:", err.Error())
		log.Printf("Failed to publish scheduled transaction event: %v", err)
		return
	}

	// Update status
	st.Status = "completed"
	println("✅ Planlanan transaction tamamlandı, ID:", st.ID)
	log.Printf("Scheduled transaction %s completed", st.ID)
}

// GetScheduledTransactions returns all scheduled transactions
func (s *Scheduler) GetScheduledTransactions() []ScheduledTransaction {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var transactions []ScheduledTransaction
	for id := range s.entries {
		// In a real implementation, you would fetch from database
		// For now, we'll return an empty slice
		transactions = append(transactions, ScheduledTransaction{ID: id})
	}

	return transactions
}

// UpdateSchedule updates the schedule of a transaction
func (s *Scheduler) UpdateSchedule(transactionID, newSchedule string) error {
	// Unschedule the old one
	if err := s.UnscheduleTransaction(transactionID); err != nil {
		return err
	}

	// Get the transaction details (in real implementation, fetch from DB)
	// For now, we'll create a dummy one
	st := &ScheduledTransaction{
		ID:       transactionID,
		Schedule: newSchedule,
	}

	// Schedule with new time
	return s.ScheduleTransaction(st)
}
