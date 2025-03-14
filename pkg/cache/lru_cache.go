package cache

import (
	"errors"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/rk-the-dev/golib-core/pkg/logger"
)

// LRUCache implements an in-memory LRU cache with TTL
type LRUCache struct {
	cache     *lru.Cache
	ttl       map[string]time.Time
	mutex     sync.Mutex
	expiryDur time.Duration
}

// NewLRUCache initializes an LRU cache with a given size and TTL
func NewLRUCache(size int, ttl time.Duration) (*LRUCache, error) {
	cache, err := lru.New(size)
	if err != nil {
		logger.Error("Failed to initialize LRU cache", map[string]interface{}{"error": err})
		return nil, err
	}

	logger.Info("LRU cache initialized", map[string]interface{}{"size": size, "ttl": ttl})
	return &LRUCache{
		cache:     cache,
		ttl:       make(map[string]time.Time),
		expiryDur: ttl,
	}, nil
}

// Set stores a value in the cache with expiration
func (c *LRUCache) Set(key string, value interface{}, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache.Add(key, value)
	c.ttl[key] = time.Now().Add(expiration)
	logger.Debug("Cache set in LRU", map[string]interface{}{"key": key, "expiration": expiration})
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expiryTime, exists := c.ttl[key]
	if !exists {
		logger.Warn("Cache miss in LRU", map[string]interface{}{"key": key})
		return nil, errors.New("key not found")
	}

	if time.Now().After(expiryTime) {
		c.cache.Remove(key)
		delete(c.ttl, key)
		logger.Warn("Cache expired in LRU", map[string]interface{}{"key": key})
		return nil, errors.New("cache expired")
	}

	value, ok := c.cache.Get(key)
	if !ok {
		logger.Warn("Cache miss in LRU", map[string]interface{}{"key": key})
		return nil, errors.New("key not found")
	}

	logger.Debug("Cache retrieved from LRU", map[string]interface{}{"key": key})
	return value, nil
}

// Delete removes a key from the cache
func (c *LRUCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache.Remove(key)
	delete(c.ttl, key)
	logger.Debug("Cache deleted from LRU", map[string]interface{}{"key": key})
}
