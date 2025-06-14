package engine

import (
	"context"
	"testing"
)

func TestMemoryEngine_Set(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "valid key-value",
			key:   "test-key",
			value: "test-value",
		},
		{
			name:  "empty value",
			key:   "test-key",
			value: "",
		},
		{
			name:  "special characters in key",
			key:   "test-key_123",
			value: "test-value",
		},
		{
			name:  "special characters in value",
			key:   "test-key",
			value: "test-value-123",
		},
	}

	for _, tt := range tests {
		tt := tt // Create new variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Enable parallel execution
			e := NewMemoryEngine()
			e.Set(context.Background(), tt.key, tt.value)

			// Verify the value was set
			got, exists := e.Get(context.Background(), tt.key)
			if !exists {
				t.Errorf("Get() key not found after Set()")
				return
			}
			if got != tt.value {
				t.Errorf("Get() = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestMemoryEngine_Get(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		setup      func(*MemoryEngine)
		want       string
		wantExists bool
	}{
		{
			name: "existing key",
			key:  "test-key",
			setup: func(e *MemoryEngine) {
				e.Set(context.Background(), "test-key", "test-value")
			},
			want:       "test-value",
			wantExists: true,
		},
		{
			name: "non-existing key",
			key:  "non-existing",
			setup: func(e *MemoryEngine) {
				// No setup needed
			},
			want:       "",
			wantExists: false,
		},
		{
			name: "empty key",
			key:  "",
			setup: func(e *MemoryEngine) {
				// No setup needed
			},
			want:       "",
			wantExists: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Create new variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Enable parallel execution
			e := NewMemoryEngine()
			if tt.setup != nil {
				tt.setup(e)
			}

			got, exists := e.Get(context.Background(), tt.key)

			if exists != tt.wantExists {
				t.Errorf("Get() exists = %v, want %v", exists, tt.wantExists)
				return
			}

			if exists && got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryEngine_Delete(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		setup func(*MemoryEngine)
	}{
		{
			name: "existing key",
			key:  "test-key",
			setup: func(e *MemoryEngine) {
				e.Set(context.Background(), "test-key", "test-value")
			},
		},
		{
			name: "non-existing key",
			key:  "non-existing",
			setup: func(e *MemoryEngine) {
				// No setup needed
			},
		},
		{
			name: "empty key",
			key:  "",
			setup: func(e *MemoryEngine) {
				// No setup needed
			},
		},
	}

	for _, tt := range tests {
		tt := tt // Create new variable for parallel tests
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Enable parallel execution
			e := NewMemoryEngine()
			if tt.setup != nil {
				tt.setup(e)
			}

			e.Delete(context.Background(), tt.key)

			// Verify the key was deleted
			_, exists := e.Get(context.Background(), tt.key)
			if exists {
				t.Errorf("Get() after Delete() key still exists")
			}
		})
	}
}

func TestMemoryEngine_Concurrent(t *testing.T) {
	t.Parallel() // Enable parallel execution for the entire concurrent test suite

	key := "test-key"
	value := "test-value"

	// Test concurrent Set operations
	t.Run("concurrent Set", func(t *testing.T) {
		t.Parallel() // Enable parallel execution for this subtest
		e := NewMemoryEngine()
		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				e.Set(context.Background(), key, value)
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		got, exists := e.Get(context.Background(), key)
		if !exists {
			t.Errorf("Get() key not found after concurrent Set()")
		}
		if got != value {
			t.Errorf("Get() = %v, want %v", got, value)
		}
	})

	// Test concurrent Get operations
	t.Run("concurrent Get", func(t *testing.T) {
		t.Parallel() // Enable parallel execution for this subtest
		e := NewMemoryEngine()
		e.Set(context.Background(), key, value)

		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				got, exists := e.Get(context.Background(), key)
				if !exists {
					t.Errorf("Get() key not found")
				}
				if got != value {
					t.Errorf("Get() = %v, want %v", got, value)
				}
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})

	// Test concurrent Delete operations
	t.Run("concurrent Delete", func(t *testing.T) {
		t.Parallel() // Enable parallel execution for this subtest
		e := NewMemoryEngine()
		e.Set(context.Background(), key, value)

		done := make(chan bool)
		for i := 0; i < 10; i++ {
			go func() {
				e.Delete(context.Background(), key)
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}

		_, exists := e.Get(context.Background(), key)
		if exists {
			t.Errorf("Get() key still exists after concurrent Delete()")
		}
	})
}
