#!/bin/bash

set -e

# Install dependencies if needed
go mod tidy

# Run all tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Display coverage
go tool cover -func=coverage.out

echo "All tests completed successfully!" 