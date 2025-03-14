package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownHelper defines the interface for managing graceful shutdown
type ShutdownHelper interface {
	RegisterShutdownHook(name string, cleanupFuncs ...func(ctx context.Context))
	WaitForShutdown()
}

// shutdownHelper implements ShutdownHelper
type shutdownHelper struct {
	hooks map[string][]func(ctx context.Context)
	mu    sync.Mutex
}

var (
	instance *shutdownHelper
	once     sync.Once
)

// NewShutdownHelper initializes and returns a ShutdownHelper instance
func NewShutdownHelper() ShutdownHelper {
	once.Do(func() {
		instance = &shutdownHelper{
			hooks: make(map[string][]func(ctx context.Context)),
		}
	})
	return instance
}

// RegisterShutdownHook adds multiple cleanup functions to be executed on shutdown
func (s *shutdownHelper) RegisterShutdownHook(name string, cleanupFuncs ...func(ctx context.Context)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hooks[name] = append(s.hooks[name], cleanupFuncs...)
}

// WaitForShutdown listens for termination signals and executes registered cleanup functions in order
func (s *shutdownHelper) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // Wait for signal
	fmt.Println("\n🔻 Graceful shutdown initiated...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, funcs := range s.hooks {
		fmt.Println("🛑 Running cleanup for:", name)
		for _, cleanup := range funcs {
			cleanup(ctx)
		}
	}
	fmt.Println("✅ Shutdown complete.")
}
