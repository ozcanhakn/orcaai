package orchestrator

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheResult represents a cached AI response
type CacheResult struct {
	Response  interface{} `json:"response"`
	CreatedAt time.Time   `json:"created_at"`
	Provider  string      `json:"provider"`
	Model     string      `json:"model"`
	CachedKey string      `json:"cached_key"`
}

// Cache interface for different cache implementations
type Cache interface {
	Get(ctx context.Context, key string) (*CacheResult, error)
	Set(ctx context.Context, key string, result *CacheResult, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache instance
func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{client: client}, nil
}

// Get retrieves a cached result by key
func (r *RedisCache) Get(ctx context.Context, key string) (*CacheResult, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, fmt.Errorf("failed to get cache entry: %w", err)
	}

	var result CacheResult
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	return &result, nil
}

// Set stores a result in cache with expiration
func (r *RedisCache) Set(ctx context.Context, key string, result *CacheResult, expiration time.Duration) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	err = r.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set cache entry: %w", err)
	}

	return nil
}

// Delete removes a cached entry
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cache entry: %w", err)
	}

	return nil
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// GenerateCacheKey creates a cache key from prompt and parameters
func GenerateCacheKey(prompt string, taskType string, provider string, model string) string {
	key := fmt.Sprintf("%s:%s:%s:%s", prompt, taskType, provider, model)
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// InMemoryCache implements Cache interface using in-memory storage
type InMemoryCache struct {
	data map[string]*CacheResult
}

// NewInMemoryCache creates a new in-memory cache instance
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]*CacheResult),
	}
}

// Get retrieves a cached result by key
func (m *InMemoryCache) Get(ctx context.Context, key string) (*CacheResult, error) {
	if result, exists := m.data[key]; exists {
		// Check if expired
		if time.Since(result.CreatedAt) < 1*time.Hour { // Default 1 hour expiration
			return result, nil
		}
		// Expired, remove it
		delete(m.data, key)
	}
	return nil, nil // Cache miss
}

// Set stores a result in cache
func (m *InMemoryCache) Set(ctx context.Context, key string, result *CacheResult, expiration time.Duration) error {
	result.CreatedAt = time.Now()
	m.data[key] = result
	return nil
}

// Delete removes a cached entry
func (m *InMemoryCache) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// Close does nothing for in-memory cache
func (m *InMemoryCache) Close() error {
	return nil
}
