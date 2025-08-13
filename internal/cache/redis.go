package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements caching using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr, password string, db int) *RedisCache {
	println("🔴 Redis cache başlatılıyor, adres:", addr)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	println("✅ Redis cache başlatıldı")
	return &RedisCache{client: client}
}

// Set sets a key-value pair with expiration
func (rc *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	println("💾 Cache'e veri yazılıyor, key:", key)

	data, err := json.Marshal(value)
	if err != nil {
		println("❌ Cache verisi marshal edilemedi:", err.Error())
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = rc.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		println("❌ Cache'e yazma hatası:", err.Error())
	} else {
		println("✅ Cache'e veri yazıldı, key:", key)
	}
	return err
}

// Get retrieves a value by key
func (rc *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to get key: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// Delete removes a key
func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (rc *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := rc.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return result > 0, nil
}

// SetNX sets a key only if it doesn't exist
func (rc *RedisCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result, err := rc.client.SetNX(ctx, key, data, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set key: %w", err)
	}

	return result, nil
}

// Incr increments a counter
func (rc *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return rc.client.Incr(ctx, key).Result()
}

// IncrBy increments a counter by a specific amount
func (rc *RedisCache) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return rc.client.IncrBy(ctx, key, value).Result()
}

// TestConnection tests the Redis connection
func (rc *RedisCache) TestConnection(ctx context.Context) error {
	println("🔴 Redis bağlantısı test ediliyor...")

	// Simple ping test
	_, err := rc.client.Ping(ctx).Result()
	if err != nil {
		println("❌ Redis ping hatası:", err.Error())
		return fmt.Errorf("redis ping failed: %w", err)
	}

	println("✅ Redis bağlantısı başarılı")
	return nil
}

// Close closes the Redis connection
func (rc *RedisCache) Close() error {
	return rc.client.Close()
}
