package storage

import (
	"concurrency-go/internal"
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

// Engine defines the interface for storage engines
type Engine interface {
	// Set stores a value for the given key
	Set(ctx context.Context, key, value string)
	// Get retrieves a value for the given key and returns true if key exists
	Get(ctx context.Context, key string) (string, bool)
	// Delete removes the value for the given key
	Delete(ctx context.Context, key string)
}

// defaultStorage implements the StorageLayer interface and handles common logic
type defaultStorage struct {
	engine Engine
	logger *zap.Logger
}

// New creates a new storage layer with the provided engine and logger
func New(engine Engine, logger *zap.Logger) internal.StorageLayer {
	return &defaultStorage{
		engine: engine,
		logger: logger,
	}
}

// Set stores a value for the given key
func (s *defaultStorage) Set(ctx context.Context, key, value string) error {
	s.logger.Info("storage layer: Set called", zap.String("key", key), zap.String("value", value))
	if key == "" {
		return errors.New("key cannot be empty")
	}
	s.engine.Set(ctx, key, value)
	return nil
}

// Get retrieves a value for the given key
func (s *defaultStorage) Get(ctx context.Context, key string) (string, error) {
	s.logger.Info("storage layer: Get called", zap.String("key", key))
	if key == "" {
		return "", errors.New("key cannot be empty")
	}
	value, exists := s.engine.Get(ctx, key)
	if !exists {
		return "", ErrKeyNotFound
	}
	return value, nil
}

// Delete removes the value for the given key
func (s *defaultStorage) Delete(ctx context.Context, key string) error {
	s.logger.Info("storage layer: Delete called", zap.String("key", key))
	if key == "" {
		return errors.New("key cannot be empty")
	}
	s.engine.Delete(ctx, key)
	return nil
}
