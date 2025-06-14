package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"concurrency-go/internal/compute/parser"
	"concurrency-go/internal/database"
	"concurrency-go/internal/storage"
	"concurrency-go/internal/storage/engine"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize components in the correct order
	// 1. Create the storage engine
	memoryEngine := engine.NewMemoryEngine()
	logger.Info("Storage engine initialized")

	// 2. Create the storage layer with the engine and logger
	storageLayer := storage.New(memoryEngine, logger)
	logger.Info("Storage layer initialized")

	// 3. Create the compute layer
	compute := parser.New(logger)
	logger.Info("Compute layer initialized")

	// 4. Create the database with compute and storage
	db := database.New(compute, storageLayer, logger)
	logger.Info("Database initialized")

	// Create context
	ctx := context.Background()

	// Initialize reader from STDIN
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to CLI Tool")
	fmt.Println("Type 'exit' to quit")
	fmt.Println("Available commands:")
	fmt.Println("  SET <key> <value>")
	fmt.Println("  GET <key>")
	fmt.Println("  DEL <key>")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Failed to read input", zap.Error(err))
			continue
		}

		// Trim whitespace and newlines
		query := strings.TrimSpace(input)

		// Check for exit command
		if query == "exit" {
			logger.Info("Exiting application")
			fmt.Println("Goodbye!")
			break
		}

		// Process the query
		cmd, err := db.HandleQuery(ctx, query)
		if err != nil {
			logger.Error("Failed to process query",
				zap.String("query", query),
				zap.Error(err))
			fmt.Printf("Error: %v\n", err)
			continue
		}

		logger.Info("Query processed successfully",
			zap.String("type", cmd.Type),
			zap.String("key", cmd.Key))

		// Print the result based on command type
		switch cmd.Type {
		case "GET":
			fmt.Println(cmd.Value)
		case "SET", "DEL":
			fmt.Println("OK")
		}
	}
}
