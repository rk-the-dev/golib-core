package main

import (
	"fmt"
	"time"

	"github.com/rk-the-dev/golib-core/pkg/cache"
	"github.com/rk-the-dev/golib-core/pkg/logger"
)

func main() {
	logger.InitializeLogger("debug", "app.log", 5, 3, 7)
	redisCache, err := cache.NewRedisCache("localhost:6379", "", 0)
	if err != nil {
		panic(err)
	}

	// Set a value in cache
	redisCache.Set("user:123", "John Doe", 10*time.Minute)

	// Retrieve the value
	value, _ := redisCache.Get("user:123")
	fmt.Println("Cached Value:", value)

	// Delete the cache
	redisCache.Delete("user:123")
}
