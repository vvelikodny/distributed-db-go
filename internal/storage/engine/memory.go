package engine

import (
	"context"
	"sync"
)

// MemoryEngine implements the Engine interface using in-memory storage
type MemoryEngine struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewMemoryEngine creates a new memory storage engine
func NewMemoryEngine() *MemoryEngine {
	return &MemoryEngine{
		store: make(map[string]string),
	}
}

// Set implements the Engine interface
func (e *MemoryEngine) Set(ctx context.Context, key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if key == "" {
		return
	}

	e.store[key] = value
}

// Get implements the Engine interface
func (e *MemoryEngine) Get(ctx context.Context, key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	value, exists := e.store[key]
	return value, exists
}

// Delete implements the Engine interface
func (e *MemoryEngine) Delete(ctx context.Context, key string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	delete(e.store, key)
}
