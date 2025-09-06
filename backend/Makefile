.PHONY: test test-orchestrator build clean

# Run all tests
test:
	go test -v ./...

# Run orchestrator tests only
test-orchestrator:
	go test -v ./orchestrator

# Run tests with coverage
test-cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run race detector
test-race:
	go test -race ./...

# Build the backend
build:
	go build -o bin/orcaai-backend .

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
install:
	go mod tidy

# Run linter
lint:
	golangci-lint run

# Run all checks
check: lint test

# Run all checks with coverage
check-cov: lint test-cov