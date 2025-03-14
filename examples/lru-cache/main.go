package main

import (
	"fmt"
	"time"

	"github.com/rk-the-dev/golib-core/pkg/cache"
	"github.com/rk-the-dev/golib-core/pkg/logger"
)

func main() {
	// ✅ Initialize the logger before using any other package
	logger.InitializeLogger("debug", "app.log", 5, 3, 7)

	// ✅ Now create the LRU cache after logger is ready
	lruCache, err := cache.NewLRUCache(100, 5*time.Minute)
	if err != nil {
		fmt.Println("Error initializing cache:", err)
		return
	}

	// Set a value in cache
	lruCache.Set("session:456", "Active", 5*time.Minute)

	// Retrieve the value
	value, err := lruCache.Get("session:456")
	if err != nil {
		fmt.Println("Cache error:", err)
	} else {
		fmt.Println("Cached Value:", value)
	}

	// Delete the cache
	lruCache.Delete("session:456")
}
