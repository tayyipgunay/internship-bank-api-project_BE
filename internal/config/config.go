package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	TokenTTL      string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
}

func LoadConfig() *Config {
	// .env yoksa durma; ortam değişkenlerinden devam et
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️  .env bulunamadı, ortam değişkenleri kullanılacak")
	}

	config := &Config{
		AppPort:       getEnvWithDefault("APP_PORT", "8080"),
		DBHost:        getEnvWithDefault("DB_HOST", "localhost"),
		DBPort:        getEnvWithDefault("DB_PORT", "5432"),
		DBUser:        getEnvWithDefault("DB_USER", "postgres"),
		DBPassword:    getEnvWithDefault("DB_PASSWORD", ""),
		DBName:        getEnvWithDefault("DB_NAME", "bankapi"),
		JWTSecret:     getEnvWithDefault("JWT_SECRET", ""),
		TokenTTL:      getEnvWithDefault("TOKEN_TTL", "1h"),
		RedisHost:     getEnvWithDefault("REDIS_HOST", "localhost"),
		RedisPort:     getEnvWithDefault("REDIS_PORT", "6379"),
		RedisPassword: getEnvWithDefault("REDIS_PASSWORD", ""),
		RedisDB:       getEnvWithDefault("REDIS_DB", "0"),
	}

	// Validate critical configurations
	if err := config.Validate(); err != nil {
		log.Printf("❌ Konfigürasyon hatası: %v", err)
	}

	return config
}

// getEnvWithDefault gets environment variable with default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Validate checks if critical configurations are set
func (c *Config) Validate() error {
	println("🔧 Konfigürasyon doğrulanıyor...")

	if c.JWTSecret == "" {
		println("⚠️ JWT_SECRET ayarlanmamış, güvenlik riski!")
	}

	if c.DBPassword == "" {
		println("⚠️ DB_PASSWORD ayarlanmamış, veritabanı bağlantısı başarısız olabilir!")
	}

	if c.RedisPassword == "" {
		println("ℹ️ REDIS_PASSWORD ayarlanmamış, Redis şifresiz çalışacak")
	}

	// Validate port numbers
	if _, err := strconv.Atoi(c.DBPort); err != nil {
		println("⚠️ DB_PORT geçersiz:", c.DBPort)
	}

	if _, err := strconv.Atoi(c.RedisPort); err != nil {
		println("⚠️ REDIS_PORT geçersiz:", c.RedisPort)
	}

	if _, err := strconv.Atoi(c.AppPort); err != nil {
		println("⚠️ APP_PORT geçersiz:", c.AppPort)
	}

	println("✅ Konfigürasyon doğrulandı")
	return nil
}
