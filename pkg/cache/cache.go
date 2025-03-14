package cache

import "time"

// CacheProvider defines the common caching interface
type CacheProvider interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}
