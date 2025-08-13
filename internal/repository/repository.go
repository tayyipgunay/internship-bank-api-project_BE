package repository

import (
	"bankapi/internal/balance"
	"bankapi/internal/transaction"
	"bankapi/internal/user"
)

// UserRepository defines storage operations for users
type UserRepository interface {
	Create(u *user.User) error
	FindByID(id uint) (*user.User, error)
	Update(u *user.User) error
	Delete(id uint) error
	List(limit, offset int) ([]user.User, error)
}

// TransactionRepository defines storage operations for transactions
type TransactionRepository interface {
	Create(t *transaction.Transaction) error
	FindByID(id uint) (*transaction.Transaction, error)
	ListLatest(limit int) ([]transaction.Transaction, error)
}

// BalanceRepository defines storage operations for balances
type BalanceRepository interface {
	GetOrCreate(userID uint) (balance.Balance, error)
	Save(b *balance.Balance) error
	History(userID uint, limit int) ([]balance.BalanceHistory, error)
}
