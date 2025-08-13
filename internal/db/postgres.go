package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() error {
	println("🗄️ PostgreSQL bağlantısı kuruluyor...")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	println("🔗 DSN:", dsn)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		println("❌ Veritabanına bağlanılamadı:", err.Error())
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := database.DB()
	if err != nil {
		println("❌ SQL DB alınamadı:", err.Error())
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	println("✅ PostgreSQL bağlantısı kuruldu")
	println("📊 Connection pool ayarlandı - MaxIdle: 10, MaxOpen: 100, MaxLifetime: 1h")
	DB = database
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		println("🗄️ Veritabanı bağlantısı kapatılıyor...")
		sqlDB, err := DB.DB()
		if err != nil {
			println("⚠️ SQL DB alınamadı:", err.Error())
			return
		}

		if err := sqlDB.Close(); err != nil {
			println("⚠️ Veritabanı kapatma hatası:", err.Error())
		} else {
			println("✅ Veritabanı bağlantısı kapatıldı")
		}
	}
}

// TestConnection tests the database connection
func TestConnection() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
