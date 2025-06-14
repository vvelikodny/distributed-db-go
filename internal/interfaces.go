package internal

import (
	"context"
	"errors"
)

var (
	ErrInvalidCommand = errors.New("invalid command")
	ErrInvalidSyntax  = errors.New("invalid syntax")
)

// Command represents a parsed command
type Command struct {
	Type  string
	Key   string
	Value string
}

// ComputeLayer defines the interface for the compute layer
type ComputeLayer interface {
	// ProcessCommand processes the input query and returns the parsed command
	ProcessCommand(ctx context.Context, query string) (*Command, error)
}

// StorageLayer defines the interface for the storage layer
type StorageLayer interface {
	// Set stores a value for the given key
	Set(ctx context.Context, key, value string) error
	// Get retrieves a value for the given key
	Get(ctx context.Context, key string) (string, error)
	// Delete removes the value for the given key
	Delete(ctx context.Context, key string) error
}
