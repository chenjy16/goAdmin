# Makefile for MCP Server Project

.PHONY: help test test-unit test-integration test-coverage test-race mock-gen clean build run

# Default target
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-unit     - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-race     - Run tests with race detection"
	@echo "  mock-gen      - Generate mock files"
	@echo "  clean         - Clean test cache and generated files"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"

# Test targets
test:
	go test -v ./...

test-unit:
	go test -v -short ./...

test-integration:
	go test -v -run Integration ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	go test -v -race ./...

# Mock generation
mock-gen:
	@echo "Generating mocks..."
	go generate ./...

# Clean targets
clean:
	go clean -testcache
	rm -f coverage.out coverage.html
	find . -name "*_mock.go" -delete

# Build and run
build:
	go build -o bin/admin cmd/main.go

run:
	go run cmd/main.go

# Development helpers
deps:
	go mod tidy
	go mod download

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

# Database related
db-migrate:
	@echo "Running database migrations..."
	# Add your migration commands here

# Docker targets (if needed)
docker-build:
	docker build -t mcp-server .

docker-run:
	docker run -p 8080:8080 mcp-server