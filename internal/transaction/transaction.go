package transaction

import (
	"fmt"
	"time"
)

type TransactionType string
type TransactionStatus string

const (
	TransactionTypeCredit   TransactionType = "credit"
	TransactionTypeDebit    TransactionType = "debit"
	TransactionTypeTransfer TransactionType = "transfer"

	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID           uint              `json:"id" gorm:"primaryKey"`
	FromUserID   *uint             `json:"from_user_id" gorm:"index"`
	ToUserID     *uint             `json:"to_user_id" gorm:"index"`
	AmountCents  int64             `json:"amount_cents" gorm:"not null;check:amount_cents>0"`
	Type         TransactionType   `json:"type" gorm:"size:20;not null"`
	Status       TransactionStatus `json:"status" gorm:"size:20;not null;index"`
	FailureCause string            `json:"failure_cause" gorm:"size:255"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type CreateCreditRequest struct {
	UserID      uint  `json:"user_id" binding:"required"`
	AmountCents int64 `json:"amount_cents" binding:"required,gt=0"`
}

type CreateDebitRequest struct {
	UserID      uint  `json:"user_id" binding:"required"`
	AmountCents int64 `json:"amount_cents" binding:"required,gt=0"`
}

type CreateTransferRequest struct {
	FromUserID  uint  `json:"from_user_id" binding:"required"`
	ToUserID    uint  `json:"to_user_id" binding:"required,nefield=FromUserID"`
	AmountCents int64 `json:"amount_cents" binding:"required,gt=0"`
}

// State management methods
func (t *Transaction) CanTransitionTo(newStatus TransactionStatus) bool {
	switch t.Status {
	case TransactionStatusPending:
		return newStatus == TransactionStatusCompleted || newStatus == TransactionStatusFailed
	case TransactionStatusCompleted:
		return false // Can't change completed transactions
	case TransactionStatusFailed:
		return false // Can't change failed transactions
	default:
		return false
	}
}

func (t *Transaction) TransitionTo(newStatus TransactionStatus) error {
	if !t.CanTransitionTo(newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", t.Status, newStatus)
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Transaction) MarkAsCompleted() error {
	return t.TransitionTo(TransactionStatusCompleted)
}

func (t *Transaction) MarkAsFailed(cause string) error {
	if err := t.TransitionTo(TransactionStatusFailed); err != nil {
		return err
	}
	t.FailureCause = cause
	return nil
}

func (t *Transaction) CanRollback() bool {
	return t.Status == TransactionStatusCompleted && t.Type != TransactionTypeTransfer
}

func (t *Transaction) Rollback() error {
	if !t.CanRollback() {
		return fmt.Errorf("transaction cannot be rolled back")
	}
	t.Status = TransactionStatusFailed
	t.FailureCause = "Rolled back by user"
	t.UpdatedAt = time.Now()
	return nil
}

// Validation methods
func (t *Transaction) Validate() error {
	if t.AmountCents <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if t.Type == TransactionTypeTransfer {
		if t.FromUserID == nil || t.ToUserID == nil {
			return fmt.Errorf("transfer transactions require both from and to user IDs")
		}
		if *t.FromUserID == *t.ToUserID {
			return fmt.Errorf("cannot transfer to same user")
		}
	}

	if t.Type == TransactionTypeCredit || t.Type == TransactionTypeDebit {
		if t.FromUserID != nil || t.ToUserID != nil {
			return fmt.Errorf("credit/debit transactions should not have from/to user IDs")
		}
	}

	return nil
}

// Business logic methods
func (t *Transaction) IsTransfer() bool {
	return t.Type == TransactionTypeTransfer
}

func (t *Transaction) IsCredit() bool {
	return t.Type == TransactionTypeCredit
}

func (t *Transaction) IsDebit() bool {
	return t.Type == TransactionTypeDebit
}

func (t *Transaction) GetAmountInTL() float64 {
	return float64(t.AmountCents) / 100.0
}

func (t *Transaction) SetAmountFromTL(amountTL float64) {
	t.AmountCents = int64(amountTL * 100)
}
