package audit

import (
	"bankapi/internal/db"
	"fmt"
	"time"
)

type AuditLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	EntityType string    `json:"entity_type" gorm:"size:50;index;not null"`
	EntityID   string    `json:"entity_id" gorm:"size:64;index;not null"`
	Action     string    `json:"action" gorm:"size:50;not null"`
	Details    string    `json:"details" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
}

func Log(entityType, entityID, action, details string) {
	println("📋 AUDIT LOG:", entityType, entityID, action, details)

	// Validate input parameters
	if entityType == "" {
		println("⚠️ Entity type boş, audit log kaydedilemedi")
		return
	}

	if entityID == "" {
		println("⚠️ Entity ID boş, audit log kaydedilemedi")
		return
	}

	if action == "" {
		println("⚠️ Action boş, audit log kaydedilemedi")
		return
	}

	// Create audit log entry
	auditLog := AuditLog{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		Details:    details,
		CreatedAt:  time.Now(),
	}

	// Try to save to database if available
	if db.DB != nil {
		if err := db.DB.Create(&auditLog).Error; err != nil {
			println("❌ Audit log veritabanına kaydedilemedi:", err.Error())
		} else {
			println("✅ Audit log veritabanına kaydedildi, ID:", auditLog.ID)
		}
	} else {
		println("ℹ️ Veritabanı bağlantısı yok, audit log sadece console'a yazıldı")
	}
}

// GetAuditLogs retrieves audit logs for a specific entity
func GetAuditLogs(entityType, entityID string) ([]AuditLog, error) {
	println("🔍 Audit loglar aranıyor, entity:", entityType, "ID:", entityID)

	if db.DB == nil {
		println("❌ Veritabanı bağlantısı yok")
		return nil, fmt.Errorf("database connection not available")
	}

	var logs []AuditLog
	query := db.DB.Where("entity_type = ? AND entity_id = ?", entityType, entityID).Order("created_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		println("❌ Audit loglar alınamadı:", err.Error())
		return nil, fmt.Errorf("failed to retrieve audit logs: %w", err)
	}

	println("✅", len(logs), "audit log bulundu")
	return logs, nil
}

// GetAuditLogsByType retrieves audit logs by entity type
func GetAuditLogsByType(entityType string) ([]AuditLog, error) {
	println("🔍 Audit loglar aranıyor, entity type:", entityType)

	if db.DB == nil {
		println("❌ Veritabanı bağlantısı yok")
		return nil, fmt.Errorf("database connection not available")
	}

	var logs []AuditLog
	query := db.DB.Where("entity_type = ?", entityType).Order("created_at DESC")

	if err := query.Find(&logs).Error; err != nil {
		println("❌ Audit loglar alınamadı:", err.Error())
		return nil, fmt.Errorf("failed to retrieve audit logs: %w", err)
	}

	println("✅", len(logs), "audit log bulundu")
	return logs, nil
}
