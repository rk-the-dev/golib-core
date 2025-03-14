package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rk-the-dev/golib-core/pkg/logger"
)

// RedisCache implements the CacheProvider interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache initializes a Redis cache
func NewRedisCache(redisAddr, redisPassword string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       db,
	})

	// Test Redis connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed to connect to Redis", map[string]interface{}{"error": err})
		return nil, err
	}

	logger.Info("Redis cache initialized successfully", nil)
	return &RedisCache{client: client}, nil
}

// Set stores a value in Redis with expiration
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		logger.Error("Failed to marshal cache data", map[string]interface{}{"error": err})
		return err
	}

	err = r.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		logger.Error("Failed to set cache in Redis", map[string]interface{}{"key": key, "error": err})
		return err
	}

	logger.Debug("Cache set in Redis", map[string]interface{}{"key": key})
	return nil
}

// Get retrieves a value from Redis
func (r *RedisCache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		logger.Warn("Cache miss in Redis", map[string]interface{}{"key": key})
		return nil, err
	}

	var value interface{}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		logger.Error("Failed to unmarshal cache data", map[string]interface{}{"key": key, "error": err})
		return nil, err
	}

	logger.Debug("Cache retrieved from Redis", map[string]interface{}{"key": key})
	return value, nil
}

// Delete removes a key from Redis
func (r *RedisCache) Delete(key string) error {
	ctx := context.Background()
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		logger.Error("Failed to delete cache from Redis", map[string]interface{}{"key": key, "error": err})
		return err
	}

	logger.Debug("Cache deleted from Redis", map[string]interface{}{"key": key})
	return nil
}
