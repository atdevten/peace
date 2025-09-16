.PHONY: help build run dev clean stop restart test test-coverage test-ci generate-mocks clean-mocks

# Default target
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the full stack (production)"
	@echo "  dev             - Run development environment"
	@echo "  clean           - Stop and remove all containers"
	@echo "  stop            - Stop all services"
	@echo "  restart         - Restart all services"
	@echo ""
	@echo "Testing:"
	@echo "  test            - Run all tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-ci         - Run tests for CI (with race detection)"
	@echo "  generate-mocks  - Generate mock files"
	@echo "  clean-mocks     - Clean mock files"

# Build the application
build: build-server build-websocket
build-server:
	@echo "Building HTTP server..."
	@go build -o bin/server ./cmd/server
build-websocket:
	@echo "Building WebSocket server..."
	@go build -o bin/websocket-server ./cmd/websocket-server

# Run production stack
run:
	docker-compose up -d

# Run development environment
dev: dev-server dev-websocket
dev-server:
	@echo "Starting HTTP server..."
	@go run ./cmd/server/main.go -config configs/config.env
dev-websocket:
	@echo "Starting WebSocket server..."
	@go run ./cmd/websocket-server/main.go -config configs/config.env

# Stop all services
stop:
	docker-compose down
	docker-compose -f docker-compose.dev.yml down

# Restart all services
restart:
	docker-compose restart
	docker-compose -f docker-compose.dev.yml restart

# Clean everything
clean:
	docker-compose down -v --remove-orphans
	docker-compose -f docker-compose.dev.yml down -v --remove-orphans
	docker system prune -f

# Testing and validation
test:
	@echo "Running tests..."
	@go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-ci:
	@echo "Running tests for CI (with race detection and coverage)..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out

# Mock generation
generate-mocks:
	@echo "Generating all mock files..."
	@./scripts/generate-mocks.sh

# Clean mocks
clean-mocks:
	@echo "Cleaning mock files..."
	@rm -rf testutils/mocks/repositories/*.go
	@rm -rf testutils/mocks/usecases/*.go
	@rm -rf testutils/mocks/services/*.go
	@echo "Mock files cleaned!"
