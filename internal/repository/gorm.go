package repository

import (
	"bankapi/internal/balance"
	"bankapi/internal/db"
	"bankapi/internal/transaction"
	"bankapi/internal/user"
)

type GormUserRepo struct{}

func (GormUserRepo) Create(u *user.User) error { return db.DB.Create(u).Error }
func (GormUserRepo) FindByID(id uint) (*user.User, error) {
	var u user.User
	if err := db.DB.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
func (GormUserRepo) Update(u *user.User) error { return db.DB.Save(u).Error }
func (GormUserRepo) Delete(id uint) error      { return db.DB.Delete(&user.User{}, id).Error }
func (GormUserRepo) List(limit, offset int) ([]user.User, error) {
	var users []user.User
	err := db.DB.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

type GormTransactionRepo struct{}

func (GormTransactionRepo) Create(t *transaction.Transaction) error { return db.DB.Create(t).Error }
func (GormTransactionRepo) FindByID(id uint) (*transaction.Transaction, error) {
	var t transaction.Transaction
	if err := db.DB.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}
func (GormTransactionRepo) ListLatest(limit int) ([]transaction.Transaction, error) {
	var txs []transaction.Transaction
	err := db.DB.Order("id DESC").Limit(limit).Find(&txs).Error
	return txs, err
}

type GormBalanceRepo struct{}

func (GormBalanceRepo) GetOrCreate(userID uint) (balance.Balance, error) {
	var b balance.Balance
	if err := db.DB.FirstOrCreate(&b, balance.Balance{UserID: userID}).Error; err != nil {
		return balance.Balance{}, err
	}
	return b, nil
}
func (GormBalanceRepo) Save(bal *balance.Balance) error { return db.DB.Save(bal).Error }
func (GormBalanceRepo) History(userID uint, limit int) ([]balance.BalanceHistory, error) {
	var hist []balance.BalanceHistory
	err := db.DB.Where("user_id = ?", userID).Order("id DESC").Limit(limit).Find(&hist).Error
	return hist, err
}
