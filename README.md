# Concurrency Go CLI

A command-line interface tool with a layered architecture for processing commands.

## Prerequisites

- Go 1.16 or higher

## Project Structure

```
.
├── cmd/
│   └── cli/          # Command-line interface
├── internal/
│   ├── compute/      # Compute layer
│   │   └── parser/   # Command parser
│   └── storage/      # Storage layer
│       └── engine/   # Storage engine
└── pkg/              # Public packages
```

## Building

1. Install dependencies:
```bash
make deps
```

2. Build the application:
```bash
make build
```

The binary will be created in the `bin` directory.

## Running

You can run the application in two ways:

1. Direct run:
```bash
make run
```

2. Build and run:
```bash
make build-run
```

## Available Commands

- `SET <key> <value>` - Store a value
- `GET <key>` - Retrieve a value
- `DEL <key>` - Delete a value
- `exit` - Exit the application

## Development

- Run tests:
```bash
make test
```

- Clean build files:
```bash
make clean
``` 