package balance

import (
	"fmt"
	"sync"
	"time"
)

type Balance struct {
	UserID      uint      `json:"user_id" gorm:"primaryKey"`
	AmountCents int64     `json:"amount_cents" gorm:"not null;default:0"`
	LastUpdated time.Time `json:"last_updated_at"`
}

type BalanceHistory struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index"`
	AmountCents int64     `json:"amount_cents"`
	CreatedAt   time.Time `json:"created_at"`
}

// Thread-safe balance operations
type ThreadSafeBalance struct {
	Balance
	mutex sync.RWMutex
}

// NewThreadSafeBalance creates a new thread-safe balance
func NewThreadSafeBalance(userID uint) *ThreadSafeBalance {
	println("💰 Thread-safe balance oluşturuluyor, kullanıcı ID:", userID)

	balance := &ThreadSafeBalance{
		Balance: Balance{
			UserID:      userID,
			AmountCents: 0,
			LastUpdated: time.Now(),
		},
	}

	println("✅ Thread-safe balance oluşturuldu")
	return balance
}

// GetAmount returns the current balance amount (thread-safe)
func (b *ThreadSafeBalance) GetAmount() int64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.AmountCents
}

// SetAmount sets the balance amount (thread-safe)
func (b *ThreadSafeBalance) SetAmount(amountCents int64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.AmountCents = amountCents
	b.LastUpdated = time.Now()
}

// AddAmount adds to the balance (thread-safe)
func (b *ThreadSafeBalance) AddAmount(amountCents int64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	newAmount := b.AmountCents + amountCents
	if newAmount < 0 {
		return fmt.Errorf("insufficient balance: cannot go below 0")
	}

	b.AmountCents = newAmount
	b.LastUpdated = time.Now()
	return nil
}

// SubtractAmount subtracts from the balance (thread-safe)
func (b *ThreadSafeBalance) SubtractAmount(amountCents int64) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.AmountCents < amountCents {
		return fmt.Errorf("insufficient balance: current=%d, requested=%d", b.AmountCents, amountCents)
	}

	b.AmountCents -= amountCents
	b.LastUpdated = time.Now()
	return nil
}

// TransferAmount transfers amount between balances (thread-safe)
func (b *ThreadSafeBalance) TransferAmount(target *ThreadSafeBalance, amountCents int64) error {
	if amountCents <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	// Lock both balances to prevent deadlock
	if b.UserID < target.UserID {
		b.mutex.Lock()
		target.mutex.Lock()
	} else {
		target.mutex.Lock()
		b.mutex.Lock()
	}
	defer b.mutex.Unlock()
	defer target.mutex.Unlock()

	if b.AmountCents < amountCents {
		return fmt.Errorf("insufficient balance for transfer")
	}

	b.AmountCents -= amountCents
	target.AmountCents += amountCents

	now := time.Now()
	b.LastUpdated = now
	target.LastUpdated = now

	return nil
}

// GetBalanceInTL returns the balance in Turkish Lira
func (b *ThreadSafeBalance) GetBalanceInTL() float64 {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return float64(b.AmountCents) / 100.0
}

// SetBalanceFromTL sets the balance from Turkish Lira amount
func (b *ThreadSafeBalance) SetBalanceFromTL(amountTL float64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.AmountCents = int64(amountTL * 100)
	b.LastUpdated = time.Now()
}
