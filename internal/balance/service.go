package balance

import (
	"bankapi/internal/audit"
	"bankapi/internal/db"
	"fmt"
	"sync"
	"time"
)

var mu sync.RWMutex

func GetOrCreateBalance(userID uint) (Balance, error) {
	println("💰 Bakiye alınıyor/oluşturuluyor, kullanıcı ID:", userID)

	var b Balance
	err := db.DB.First(&b, "user_id = ?", userID).Error
	if err == nil {
		println("✅ Mevcut bakiye bulundu:", b.AmountCents, "kuruş")
		return b, nil
	}

	println("🆕 Yeni bakiye oluşturuluyor...")
	b = Balance{UserID: userID, AmountCents: 0, LastUpdated: time.Now()}
	if err := db.DB.Create(&b).Error; err != nil {
		println("❌ Bakiye oluşturulamadı:", err.Error())
		return Balance{}, fmt.Errorf("failed to create balance: %w", err)
	}

	println("✅ Yeni bakiye oluşturuldu")
	return b, nil
}

func Credit(userID uint, amount int64) error {
	println("💳 Kredi işlemi, kullanıcı ID:", userID, "miktar:", amount, "kuruş")

	if amount <= 0 {
		println("❌ Geçersiz kredi miktarı:", amount)
		return fmt.Errorf("credit amount must be positive")
	}

	mu.Lock()
	defer mu.Unlock()

	var b Balance
	if err := db.DB.FirstOrCreate(&b, Balance{UserID: userID}).Error; err != nil {
		println("❌ Bakiye alınamadı/oluşturulamadı:", err.Error())
		return fmt.Errorf("failed to get/create balance: %w", err)
	}

	oldAmount := b.AmountCents
	b.AmountCents += amount
	b.LastUpdated = time.Now()

	if err := db.DB.Save(&b).Error; err != nil {
		println("❌ Bakiye güncellenemedi:", err.Error())
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Create balance history
	if err := db.DB.Create(&BalanceHistory{UserID: userID, AmountCents: b.AmountCents}).Error; err != nil {
		println("⚠️ Bakiye geçmişi oluşturulamadı:", err.Error())
	}

	// Audit log
	audit.Log("balance", fmt.Sprintf("%d", userID), "credit", fmt.Sprintf("+%d -> %d", amount, b.AmountCents))

	println("✅ Kredi işlemi başarılı:", oldAmount, "->", b.AmountCents, "kuruş")
	return nil
}

func Debit(userID uint, amount int64) error {
	println("💸 Debit işlemi, kullanıcı ID:", userID, "miktar:", amount, "kuruş")

	if amount <= 0 {
		println("❌ Geçersiz debit miktarı:", amount)
		return fmt.Errorf("debit amount must be positive")
	}

	mu.Lock()
	defer mu.Unlock()

	var b Balance
	if err := db.DB.FirstOrCreate(&b, Balance{UserID: userID}).Error; err != nil {
		println("❌ Bakiye alınamadı/oluşturulamadı:", err.Error())
		return fmt.Errorf("failed to get/create balance: %w", err)
	}

	if b.AmountCents < amount {
		println("❌ Yetersiz bakiye:", b.AmountCents, "<", amount)
		return ErrInsufficientFunds
	}

	oldAmount := b.AmountCents
	b.AmountCents -= amount
	b.LastUpdated = time.Now()

	if err := db.DB.Save(&b).Error; err != nil {
		println("❌ Bakiye güncellenemedi:", err.Error())
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// Create balance history
	if err := db.DB.Create(&BalanceHistory{UserID: userID, AmountCents: b.AmountCents}).Error; err != nil {
		println("⚠️ Bakiye geçmişi oluşturulamadı:", err.Error())
	}

	// Audit log
	audit.Log("balance", fmt.Sprintf("%d", userID), "debit", fmt.Sprintf("-%d -> %d", amount, b.AmountCents))

	println("✅ Debit işlemi başarılı:", oldAmount, "->", b.AmountCents, "kuruş")
	return nil
}

var ErrInsufficientFunds = &insufficientFundsError{}

type insufficientFundsError struct{}

func (insufficientFundsError) Error() string { return "yetersiz bakiye" }
