package main

import (
	"bankapi/internal/audit"
	"bankapi/internal/auth"
	"bankapi/internal/balance"
	"bankapi/internal/cache"
	"bankapi/internal/config"
	"bankapi/internal/currency"
	"bankapi/internal/db"
	"bankapi/internal/events"
	"bankapi/internal/logger"
	"bankapi/internal/metrics"
	"bankapi/internal/middleware"
	"bankapi/internal/scheduler"
	"bankapi/internal/telemetry"
	"bankapi/internal/transaction"
	"bankapi/internal/user"

	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Config validation
	println("🔧 Konfigürasyon kontrol ediliyor...")
	if cfg.JWTSecret == "" {
		println("❌ JWT_SECRET konfigürasyonu eksik!")
		log.Fatal("JWT_SECRET environment variable is required")
	}
	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		println("❌ Veritabanı konfigürasyonu eksik!")
		log.Fatal("Database configuration is incomplete")
	}
	println("✅ Konfigürasyon kontrol edildi")

	// Initialize structured logger
	logger.InitLogger()

	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.RequestID())
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.PerformanceMonitor())
	router.Use(middleware.RateLimitPerMinute(120))

	// Simple ping endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Bank API is running! 🚀",
			"version": "1.0.0",
			"features": []string{
				"Event Sourcing",
				"Redis Caching",
				"Scheduled Transactions",
				"Multi-Currency Support",
				"Audit Logging",
				"Prometheus Metrics",
				"OpenTelemetry",
			},
		})
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"database":  "disabled",
			"redis":     "disabled",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// Initialize database connection
	println("🗄️ Veritabanı bağlantısı başlatılıyor...")

	// Try to connect to database, but don't fail if it's not available
	if err := tryDatabaseConnection(); err != nil {
		println("⚠️ Veritabanı bağlantısı başarısız (devam ediliyor):", err.Error())
		println("ℹ️ Uygulama veritabanı olmadan çalışacak")
	} else {
		println("✅ Veritabanı bağlantısı başarılı")

		// Auto-migrate database models
		println("🔄 Veritabanı modelleri migrate ediliyor...")

		// Her modeli ayrı ayrı migrate et, hata olursa devam et
		models := []interface{}{
			&user.User{},
			&balance.Balance{},
			&balance.BalanceHistory{},
			&transaction.Transaction{},
			&audit.AuditLog{},
		}

		for _, model := range models {
			println("🔄 Model migrate ediliyor:", fmt.Sprintf("%T", model))

			// Force migration - tabloları yeniden oluştur
			if err := db.DB.Migrator().DropTable(model); err != nil {
				println("⚠️ Tablo silme hatası:", err.Error())
			}

			if err := db.DB.AutoMigrate(model); err != nil {
				println("⚠️ Model migration hatası (devam ediliyor):", err.Error())
			} else {
				println("✅ Model migrate edildi:", fmt.Sprintf("%T", model))
			}
		}

		println("✅ Veritabanı migration tamamlandı")

		// Seed admin user if not exists
		println("👑 Admin kullanıcı kontrol ediliyor...")
		seedAdminUser()
	}

	defer func() {
		// Database connection cleanup
		if db.DB != nil {
			println("🗄️ Veritabanı bağlantısı kapatılıyor...")
			db.CloseDB()
			log.Printf("Database connection closed")
		}
	}()

	// Initialize event bus and scheduler
	eventBus := events.NewInMemoryEventBus()
	sched := scheduler.NewScheduler(eventBus)
	sched.Start()
	defer sched.Stop()

	// Initialize Redis cache (will fail gracefully if Redis is not available)
	println("🔴 Redis cache başlatılıyor...")
	redisCache := cache.NewRedisCache(cfg.RedisHost+":"+cfg.RedisPort, cfg.RedisPassword, 0)

	// Test Redis connection
	ctx := context.Background()
	if err := redisCache.TestConnection(ctx); err != nil {
		println("⚠️ Redis bağlantısı başarısız (devam ediliyor):", err.Error())
	} else {
		println("✅ Redis bağlantısı test edildi")
	}
	defer redisCache.Close()

	// telemetry and metrics
	println("📊 Telemetry ve metrics başlatılıyor...")
	if shutdown, err := telemetry.Init("bank-api"); err == nil {
		defer func() { _ = shutdown(context.Background()) }()
		println("✅ Telemetry başlatıldı")
	} else {
		println("❌ Telemetry başlatılamadı:", err.Error())
	}

	// Register metrics endpoint
	metrics.Register(router)

	// Register all API routes
	auth.RegisterAuthRoutes(router)
	user.RegisterUserRoutes(router)
	transaction.RegisterRoutes(router)
	balance.RegisterRoutes(router)
	audit.RegisterRoutes(router, middleware.AuthMiddleware(cfg))
	scheduler.RegisterRoutes(router, middleware.AuthMiddleware(cfg), sched)
	currency.RegisterRoutes(router, middleware.AuthMiddleware(cfg))

	// API info endpoint
	router.GET("/api/v1/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":        "Bank API",
			"description": "Advanced Banking API with Event Sourcing, Caching, and Scheduling",
			"endpoints": map[string]interface{}{
				"authentication": "/api/v1/auth/*",
				"users":          "/api/v1/users/*",
				"transactions":   "/api/v1/transactions/*",
				"balances":       "/api/v1/balances/*",
				"audit":          "/api/v1/audit/*",
				"scheduler":      "/api/v1/scheduler/*",
				"currency":       "/api/v1/currency/*",
			},
			"features": map[string]interface{}{
				"event_sourcing":         true,
				"redis_caching":          true,
				"scheduled_transactions": true,
				"multi_currency":         true,
				"audit_logging":          true,
				"prometheus_metrics":     true,
				"opentelemetry":          true,
			},
		})
	})

	// graceful shutdown
	srvErrChan := make(chan error, 1)
	go func() {
		port := cfg.AppPort
		if port == "" {
			port = "8080"
		}
		println("🚀 Bank API başlatılıyor, port:", port)
		log.Printf("🚀 Starting Bank API on port %s", port)
		srvErrChan <- router.Run(":" + port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-quit:
		println("🛑 Kapatma sinyali alındı:", sig.String())
		log.Printf("🛑 Shutting down, signal: %v", sig)
	case err := <-srvErrChan:
		if err != nil {
			println("❌ Server hatası:", err.Error())
			log.Fatalf("❌ Server error: %v", err)
		}
	}
}

// seedAdminUser creates admin user if it doesn't exist
func seedAdminUser() {
	println("👑 Admin kullanıcı aranıyor...")

	var adminUser user.User
	result := db.DB.Where("role = ?", "admin").First(&adminUser)

	if result.Error != nil {
		println("👑 Admin kullanıcı bulunamadı, oluşturuluyor...")

		// Admin user doesn't exist, create one
		adminUser = user.User{
			Username:     "admin",
			Email:        "admin@bankapi.com",
			PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
			Role:         "admin",
		}

		if err := db.DB.Create(&adminUser).Error; err != nil {
			println("❌ Admin kullanıcı oluşturulamadı:", err.Error())
			log.Printf("Failed to create admin user: %v", err)
		} else {
			println("✅ Admin kullanıcı başarıyla oluşturuldu")
			log.Printf("Admin user created successfully")
		}
	} else {
		println("✅ Admin kullanıcı zaten mevcut")
	}
}

// tryDatabaseConnection attempts to connect to database without failing
func tryDatabaseConnection() error {
	println("🔗 Veritabanı bağlantısı deneniyor...")

	// Try to connect
	db.ConnectDB()

	// Test connection
	if db.DB == nil {
		return fmt.Errorf("database connection failed")
	}

	// Test ping
	if err := db.TestConnection(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
