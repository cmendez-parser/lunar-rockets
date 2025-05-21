# Lunar Rockets

A Go microservice that consumes and processes rocket state messages and exposes rocket information through a REST API.

## Lite Clean Architecture

This project follows clean architecture principles with clear separation of concerns:

- **Domain**: Core business entities and repository interfaces
- **Repository**: Data access layer
- **Usecase**: Business logic
- **Http**: HTTP Rest controllers and routing

## Features

- Receive and process rocket state messages events.
- Handle out-of-order and duplicate messages.
- Store rocket state in SQLite database.
- Expose REST API for querying rocket information.

## API Endpoints

The API documentation is available through Swagger UI at `http://localhost:8088/swagger/index.html` when the service is running.

Available endpoints:
- `POST /messages`: Receive rocket messages via webhook
- `GET /rockets`: List all rockets with optional sorting
- `GET /rockets/{channel}`: Get a specific rocket by channel ID 

## Requirements

- Go 1.24 or higher
- SQLite3

## Running the Service

If your system is darwin_arm64:
```bash
# Run the service
./lunar-rockets   
```

If your system is different than darwin_arm64:
```bash
# Build the service
go build -o lunar-rockets ./cmd

# Run the service
./lunar-rockets   
```

## Running the Test Program

Use the provided test program to simulate rocket messages:

```bash
./rockets launch "http://localhost:8088/messages" --message-delay=500ms --concurrency-level=1
```

## Environment Variables

- `SERVER_ADDRESS`: HTTP server address (default: ":8088")
- `DB_PATH`: Path to SQLite database (default: "data/rockets.db")

## Project Structure

```
.
├── cmd/               # Application entry point
├── configs/           # Configuration files
├── data/              # Data storage directory
├── db/                # Database connection and migrations
├── domain/            # Domain models and interfaces
├── http/              # HTTP controllers and routing
├── repository/        # Data access implementations
├── test/              # Test utilities and mocks
└── usecase/           # Business logic implementations
```

## Testing

The project uses table-driven tests with parallel execution for unit testing. Mock implementations are provided for all interfaces to support testing without external dependencies.

### Running Tests

```bash
# Run all tests
./test/scripts/run_tests.sh

# Run tests for a specific package
go test -v ./usecase/...
```

### Test Structure

- `test/mocks`: Contains mock implementations of interfaces
- `test/helper`: Contains helper functions for testing

### Test Coverage

Test coverage is reported when running the test script. The goal is to maintain high test coverage (>80%) for all business logic. 