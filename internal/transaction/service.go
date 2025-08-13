package transaction

import (
	"bankapi/internal/audit"
	"bankapi/internal/balance"
	"bankapi/internal/db"
	"fmt"
)

func ApplyCredit(userID uint, amount int64) (*Transaction, error) {
	println("💳 Kredi işlemi uygulanıyor, kullanıcı ID:", userID, "miktar:", amount, "kuruş")

	if amount <= 0 {
		println("❌ Geçersiz kredi miktarı:", amount)
		return nil, fmt.Errorf("credit amount must be positive")
	}

	tx := &Transaction{ToUserID: &userID, AmountCents: amount, Type: TransactionTypeCredit, Status: TransactionStatusPending}

	if err := db.DB.Create(tx).Error; err != nil {
		println("❌ Transaction oluşturulamadı:", err.Error())
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	println("📝 Transaction oluşturuldu, ID:", tx.ID)

	if err := balance.Credit(userID, amount); err != nil {
		println("❌ Bakiye kredisi başarısız:", err.Error())
		tx.Status = TransactionStatusFailed
		tx.FailureCause = err.Error()
		if saveErr := db.DB.Save(tx).Error; saveErr != nil {
			println("⚠️ Failed transaction kaydedilemedi:", saveErr.Error())
		}
		return tx, err
	}

	tx.Status = TransactionStatusCompleted
	if err := db.DB.Save(tx).Error; err != nil {
		println("❌ Transaction güncellenemedi:", err.Error())
		return tx, fmt.Errorf("failed to update transaction: %w", err)
	}

	audit.Log("transaction", fmt.Sprintf("%d", tx.ID), "credit", fmt.Sprintf("to=%d amount=%d", userID, amount))
	println("✅ Kredi işlemi başarıyla tamamlandı, transaction ID:", tx.ID)
	return tx, nil
}

func ApplyDebit(userID uint, amount int64) (*Transaction, error) {
	println("💸 Debit işlemi uygulanıyor, kullanıcı ID:", userID, "miktar:", amount, "kuruş")

	if amount <= 0 {
		println("❌ Geçersiz debit miktarı:", amount)
		return nil, fmt.Errorf("debit amount must be positive")
	}

	tx := &Transaction{FromUserID: &userID, AmountCents: amount, Type: TransactionTypeDebit, Status: TransactionStatusPending}

	if err := db.DB.Create(tx).Error; err != nil {
		println("❌ Transaction oluşturulamadı:", err.Error())
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	println("📝 Transaction oluşturuldu, ID:", tx.ID)

	if err := balance.Debit(userID, amount); err != nil {
		println("❌ Bakiye debiti başarısız:", err.Error())
		tx.Status = TransactionStatusFailed
		tx.FailureCause = err.Error()
		if saveErr := db.DB.Save(tx).Error; saveErr != nil {
			println("⚠️ Failed transaction kaydedilemedi:", saveErr.Error())
		}
		return tx, err
	}

	tx.Status = TransactionStatusCompleted
	if err := db.DB.Save(tx).Error; err != nil {
		println("❌ Transaction güncellenemedi:", err.Error())
		return tx, fmt.Errorf("failed to update transaction: %w", err)
	}

	audit.Log("transaction", fmt.Sprintf("%d", tx.ID), "debit", fmt.Sprintf("from=%d amount=%d", userID, amount))
	println("✅ Debit işlemi başarıyla tamamlandı, transaction ID:", tx.ID)
	return tx, nil
}

func ApplyTransfer(fromID, toID uint, amount int64) (*Transaction, error) {
	println("🔄 Transfer işlemi uygulanıyor, from:", fromID, "to:", toID, "miktar:", amount, "kuruş")

	if amount <= 0 {
		println("❌ Geçersiz transfer miktarı:", amount)
		return nil, fmt.Errorf("transfer amount must be positive")
	}

	if fromID == toID {
		println("❌ Aynı kullanıcıya transfer yapılamaz")
		return nil, fmt.Errorf("cannot transfer to same user")
	}

	txModel := &Transaction{FromUserID: &fromID, ToUserID: &toID, AmountCents: amount, Type: TransactionTypeTransfer, Status: TransactionStatusPending}

	if err := db.DB.Create(txModel).Error; err != nil {
		println("❌ Transfer transaction oluşturulamadı:", err.Error())
		return nil, fmt.Errorf("failed to create transfer transaction: %w", err)
	}

	println("📝 Transfer transaction oluşturuldu, ID:", txModel.ID)

	// First debit from source account
	if err := balance.Debit(fromID, amount); err != nil {
		println("❌ Kaynak hesaptan debit başarısız:", err.Error())
		txModel.Status = TransactionStatusFailed
		txModel.FailureCause = err.Error()
		if saveErr := db.DB.Save(txModel).Error; saveErr != nil {
			println("⚠️ Failed transfer transaction kaydedilemedi:", saveErr.Error())
		}
		return txModel, err
	}

	println("✅ Kaynak hesaptan debit başarılı")

	// Then credit to destination account
	if err := balance.Credit(toID, amount); err != nil {
		println("❌ Hedef hesaba kredi başarısız, rollback yapılıyor:", err.Error())
		// rollback debit if credit fails
		if rollbackErr := balance.Credit(fromID, amount); rollbackErr != nil {
			println("❌ Rollback başarısız:", rollbackErr.Error())
		} else {
			println("✅ Rollback başarılı")
		}

		txModel.Status = TransactionStatusFailed
		txModel.FailureCause = err.Error()
		if saveErr := db.DB.Save(txModel).Error; saveErr != nil {
			println("⚠️ Failed transfer transaction kaydedilemedi:", saveErr.Error())
		}
		return txModel, err
	}

	println("✅ Hedef hesaba kredi başarılı")

	txModel.Status = TransactionStatusCompleted
	if err := db.DB.Save(txModel).Error; err != nil {
		println("❌ Transfer transaction güncellenemedi:", err.Error())
		return txModel, fmt.Errorf("failed to update transfer transaction: %w", err)
	}

	audit.Log("transaction", fmt.Sprintf("%d", txModel.ID), "transfer", fmt.Sprintf("from=%d to=%d amount=%d", fromID, toID, amount))
	println("✅ Transfer işlemi başarıyla tamamlandı, transaction ID:", txModel.ID)
	return txModel, nil
}
