.PHONY: test test-client test-commands build clean

# Run all tests
test:
	go test -v ./...

# Run client tests only
test-client:
	go test -v ./client

# Run commands tests only
test-commands:
	go test -v ./commands

# Build the CLI
build:
	go build -o bin/orcaai .

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
install:
	go mod tidy

# Run tests with coverage
test-cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run race detector
test-race:
	go test -race ./...