package database

import (
	"context"

	"concurrency-go/internal"

	"go.uber.org/zap"
)

// Database represents the main database instance that combines compute and storage layers
type Database struct {
	compute internal.ComputeLayer
	storage internal.StorageLayer
	logger  *zap.Logger
}

// New creates a new Database instance with the provided compute and storage implementations
func New(compute internal.ComputeLayer, storage internal.StorageLayer, logger *zap.Logger) *Database {
	return &Database{
		compute: compute,
		storage: storage,
		logger:  logger,
	}
}

// HandleQuery processes a query through the compute layer and executes the command using storage
func (db *Database) HandleQuery(ctx context.Context, query string) (*internal.Command, error) {
	db.logger.Info("processing query", zap.String("query", query))

	// Parse the command
	cmd, err := db.compute.ProcessCommand(ctx, query)
	if err != nil {
		db.logger.Error("failed to process command",
			zap.String("query", query),
			zap.Error(err))
		return nil, err
	}

	// Execute the command using storage
	switch cmd.Type {
	case "SET":
		err = db.storage.Set(ctx, cmd.Key, cmd.Value)
		if err != nil {
			db.logger.Error("failed to set value",
				zap.String("key", cmd.Key),
				zap.String("value", cmd.Value),
				zap.Error(err))
			return nil, err
		}
		db.logger.Info("value set successfully",
			zap.String("key", cmd.Key),
			zap.String("value", cmd.Value))
		return cmd, nil

	case "GET":
		value, err := db.storage.Get(ctx, cmd.Key)
		if err != nil {
			db.logger.Error("failed to get value",
				zap.String("key", cmd.Key),
				zap.Error(err))
			return nil, err
		}
		cmd.Value = value
		db.logger.Info("value retrieved successfully",
			zap.String("key", cmd.Key),
			zap.String("value", value))
		return cmd, nil

	case "DEL":
		err = db.storage.Delete(ctx, cmd.Key)
		if err != nil {
			db.logger.Error("failed to delete key",
				zap.String("key", cmd.Key),
				zap.Error(err))
			return nil, err
		}
		db.logger.Info("key deleted successfully",
			zap.String("key", cmd.Key))
		return cmd, nil

	default:
		db.logger.Error("invalid command type",
			zap.String("type", cmd.Type))
		return nil, internal.ErrInvalidCommand
	}
}
