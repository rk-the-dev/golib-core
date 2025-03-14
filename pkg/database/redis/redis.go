package redishelper

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient defines the interface for Redis operations
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	GetClient() *redis.Client // Provides direct access to the Redis client
	Close() error             // Optional closing of the Redis connection
}

// RedisConfig defines Redis connection configurations
type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"localhost"`
	Port     int    `env:"REDIS_PORT" envDefault:"6379"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

// redisClient is an implementation of RedisClient
type redisClient struct {
	client *redis.Client
}

var (
	instance *redisClient
	once     sync.Once
)

// NewRedisClient initializes and returns a Redis connection
func NewRedisClient(cfg *RedisConfig) (RedisClient, error) {
	var err error
	once.Do(func() {
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		})
		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err = rdb.Ping(ctx).Err(); err != nil {
			fmt.Println("❌ Failed to connect to Redis:", err)
			instance = nil
			return
		}
		fmt.Println("✅ Connected to Redis on", addr)
		instance = &redisClient{client: rdb}
	})
	return instance, err
}

// Set stores a key-value pair in Redis
func (r *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key
func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete removes a key from Redis
func (r *redisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// GetClient returns the underlying Redis client for custom operations
func (r *redisClient) GetClient() *redis.Client {
	return r.client
}

// Close closes the Redis connection
func (r *redisClient) Close() error {
	return r.client.Close()
}
