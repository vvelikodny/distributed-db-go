.PHONY: build clean test run

# Binary name
BINARY_NAME=cli

# Build the application
build:
	go build -o bin/$(BINARY_NAME) cmd/cli/main.go

# Clean build files
clean:
	go clean
	rm -f bin/$(BINARY_NAME)

# Run tests
test:
	go test ./...

# Run the application
run:
	go run cmd/cli/main.go

# Build and run
build-run: build
	./bin/$(BINARY_NAME)

# Install dependencies
deps:
	go mod tidy 