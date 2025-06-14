package storage

import (
	"context"
	"testing"

	"concurrency-go/internal"
	"concurrency-go/internal/storage/engine"

	"go.uber.org/zap"
)

func TestDefaultStorage_Set(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{
			name:    "valid key-value",
			key:     "test-key",
			value:   "test-value",
			wantErr: false,
		},
		{
			name:    "empty key",
			key:     "",
			value:   "test-value",
			wantErr: true,
		},
		{
			name:    "empty value",
			key:     "test-key",
			value:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger, _ := zap.NewDevelopment()
			storage := New(engine.NewMemoryEngine(), logger)

			err := storage.Set(context.Background(), tt.key, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the value was set
				got, err := storage.Get(context.Background(), tt.key)
				if err != nil {
					t.Errorf("Get() error = %v", err)
					return
				}
				if got != tt.value {
					t.Errorf("Get() = %v, want %v", got, tt.value)
				}
			}
		})
	}
}

func TestDefaultStorage_Get(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setup    func(internal.StorageLayer)
		want     string
		wantErr  bool
		checkErr error
	}{
		{
			name: "existing key",
			key:  "test-key",
			setup: func(s internal.StorageLayer) {
				s.Set(context.Background(), "test-key", "test-value")
			},
			want:     "test-value",
			wantErr:  false,
			checkErr: nil,
		},
		{
			name: "non-existing key",
			key:  "non-existing",
			setup: func(s internal.StorageLayer) {
				// No setup needed
			},
			want:     "",
			wantErr:  true,
			checkErr: ErrKeyNotFound,
		},
		{
			name: "empty key",
			key:  "",
			setup: func(s internal.StorageLayer) {
				// No setup needed
			},
			want:     "",
			wantErr:  true,
			checkErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger, _ := zap.NewDevelopment()
			storage := New(engine.NewMemoryEngine(), logger)

			if tt.setup != nil {
				tt.setup(storage)
			}

			got, err := storage.Get(context.Background(), tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkErr != nil && err != tt.checkErr {
				t.Errorf("Get() error = %v, want %v", err, tt.checkErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultStorage_Delete(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		setup    func(internal.StorageLayer)
		wantErr  bool
		checkErr error
	}{
		{
			name: "existing key",
			key:  "test-key",
			setup: func(s internal.StorageLayer) {
				s.Set(context.Background(), "test-key", "test-value")
			},
			wantErr:  false,
			checkErr: nil,
		},
		{
			name: "non-existing key",
			key:  "non-existing",
			setup: func(s internal.StorageLayer) {
				// No setup needed
			},
			wantErr:  false,
			checkErr: nil,
		},
		{
			name: "empty key",
			key:  "",
			setup: func(s internal.StorageLayer) {
				// No setup needed
			},
			wantErr:  true,
			checkErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			logger, _ := zap.NewDevelopment()
			storage := New(engine.NewMemoryEngine(), logger)

			if tt.setup != nil {
				tt.setup(storage)
			}

			err := storage.Delete(context.Background(), tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkErr != nil && err != tt.checkErr {
				t.Errorf("Delete() error = %v, want %v", err, tt.checkErr)
				return
			}

			if !tt.wantErr {
				// Verify the key was deleted
				_, err := storage.Get(context.Background(), tt.key)
				if err != ErrKeyNotFound {
					t.Errorf("Get() after Delete() error = %v, want %v", err, ErrKeyNotFound)
				}
			}
		})
	}
}
