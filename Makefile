.PHONY: help build run dev clean logs ps stop restart quotes-import quotes-import-fast quotes-dry-run quotes-check

# Default target
help:
	@echo "Available commands:"
	@echo "  build           - Build the application Docker image"
	@echo "  run             - Run the full stack (production)"
	@echo "  dev             - Run development environment with hot reload"
	@echo "  clean           - Stop and remove all containers and volumes"
	@echo "  stop            - Stop all services"
	@echo "  restart         - Restart all services"
	@echo "  health          - Check service health"
	@echo ""
	@echo "Quote management:"
	@echo "  quotes-import      - Import quotes from CSV to database (10 workers)"
	@echo "  quotes-import-fast - Import quotes from CSV to database (20 workers, fast mode)"
	@echo "  quotes-dry-run     - Test quote import without inserting to DB"

# Build the application
build:
	docker-compose build

# Run production stack
run:
	docker-compose up -d

# Run development environment
dev:
	docker-compose -f docker-compose.dev.yml up -d

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

# Health check
health:
	@echo "Checking service health..."
	@curl -f http://localhost:8080/health || echo "App is not healthy"
	@curl -f http://localhost:5432 || echo "PostgreSQL is not healthy"
	@curl -f http://localhost:8086/health || echo "InfluxDB is not healthy"

# Quote management commands
quotes-import:
	@echo "Building quotes importer..."
	@go build -o bin/quotes ./cmd/data-importer
	@echo "Importing quotes from CSV to database..."
	@./bin/quotes --config=configs/config.yml --workers=10 --batch-size=100

quotes-import-fast:
	@echo "Building quotes importer..."
	@go build -o bin/quotes ./cmd/data-importer
	@echo "Importing quotes from CSV to database (fast mode)..."
	@./bin/quotes --config=configs/config.yml --workers=20 --batch-size=200

quotes-dry-run:
	@echo "Building quotes importer..."
	@go build -o bin/quotes ./cmd/data-importer
	@echo "Testing quote import (dry run)..."
	@./bin/quotes --config=configs/config.yml --dry-run --workers=5 --batch-size=50
