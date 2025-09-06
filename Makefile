.PHONY: help build run dev clean logs ps stop restart quotes-import quotes-import-fast quotes-dry-run quotes-import-zip quotes-import-zip-dry data-reset data-migrate data-seed test test-coverage db-connect db-backup clean-bin clean-logs logs-app logs-db logs-redis

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
	@echo "  quotes-import-zip  - Import quotes from ZIP file containing CSV files"
	@echo "  quotes-import-zip-dry - Test ZIP import without inserting to DB"
	@echo ""
	@echo "Data management:"
	@echo "  data-reset         - Reset database (drop all data)"
	@echo "  data-migrate       - Run database migrations"
	@echo "  data-seed          - Seed database with sample data"

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
	@go run ./cmd/server/main.go -config configs/config.yml
dev-websocket:
	@echo "Starting WebSocket server..."
	@go run ./cmd/websocket-server/main.go -config configs/config.yml

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
	@curl -f http://localhost:8080/health || echo "HTTP server is not healthy"
	@curl -f http://localhost:8081/ws/health || echo "WebSocket server is not healthy"
	@curl -f http://localhost:5432 || echo "PostgreSQL is not healthy"
	@curl -f http://localhost:6380 || echo "Redis is not healthy"

# ZIP file import commands
quotes-import-zip:
	@echo "Building quotes importer..."
	@go build -o bin/quotes ./cmd/data-importer
	@echo "Importing quotes from ZIP file..."
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Please specify FILE parameter"; \
		echo "Example: make quotes-import-zip FILE=quotes.zip"; \
		exit 1; \
	fi
	@echo "Processing file: $(FILE)"
	@./bin/quotes --config=configs/config.yml --file="$(FILE)" --workers=10 --batch-size=100

quotes-import-zip-dry:
	@echo "Building quotes importer..."
	@go build -o bin/quotes ./cmd/data-importer
	@echo "Testing ZIP import (dry run)..."
	@if [ -z "$(FILE)" ]; then \
		echo "Error: Please specify FILE parameter"; \
		echo "Example: make quotes-import-zip-dry FILE=quotes.zip"; \
		exit 1; \
	fi
	@echo "Processing file: $(FILE)"
	@./bin/quotes --config=configs/config.yml --file="$(FILE)" --dry-run --workers=5 --batch-size=50
	
# Data management commands
data-reset:
	@echo "⚠️  WARNING: This will delete ALL data from the database!"
	@read -p "Are you sure? Type 'yes' to confirm: " confirm; \
	if [ "$$confirm" = "yes" ]; then \
		echo "Stopping services..."; \
		docker-compose -f docker-compose.dev.yml down -v; \
		echo "Removing database volume..."; \
		docker volume rm peace_postgres_dev_data || true; \
		echo "Starting fresh database..."; \
		docker-compose -f docker-compose.dev.yml up -d postgres; \
		echo "Waiting for database to be ready..."; \
		sleep 10; \
		echo "Running migrations..."; \
		make data-migrate; \
		echo "Database reset complete!"; \
	else \
		echo "Database reset cancelled."; \
	fi

data-migrate:
	@echo "Running database migrations..."
	@docker exec -it peace_postgres_dev psql -U postgres -d peace -c "SELECT version_id, is_applied, tstamp FROM goose_db_version ORDER BY version_id;" || \
	(echo "No migrations found, running initial setup..." && \
	docker exec -it peace_postgres_dev psql -U postgres -d peace -c "CREATE TABLE IF NOT EXISTS goose_db_version (version_id BIGINT PRIMARY KEY, is_applied BOOLEAN NOT NULL, tstamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	@echo "Migrations completed!"

data-seed:
	@echo "Seeding database with sample data..."
	@if [ -f "quotes.csv" ]; then \
		echo "Found quotes.csv, importing sample quotes..."; \
		make quotes-import; \
	else \
		echo "No quotes.csv found. Please create a CSV file with quotes data."; \
	fi

# Development utilities
logs:
	@docker-compose -f docker-compose.dev.yml logs -f

logs-app:
	@docker-compose -f docker-compose.dev.yml logs -f app

logs-db:
	@docker-compose -f docker-compose.dev.yml logs -f postgres

logs-redis:
	@docker-compose -f docker-compose.dev.yml logs -f redis

# Testing and validation
test:
	@echo "Running tests..."
	@go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Database utilities
db-connect:
	@echo "Connecting to PostgreSQL database..."
	@docker exec -it peace_postgres_dev psql -U postgres -d peace

db-backup:
	@echo "Creating database backup..."
	@docker exec peace_postgres_dev pg_dump -U postgres -d peace > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup created: backup_$(shell date +%Y%m%d_%H%M%S).sql"

# Cleanup utilities
clean-bin:
	@echo "Cleaning binary files..."
	@rm -rf bin/
	@echo "Binary files cleaned!"

clean-logs:
	@echo "Cleaning log files..."
	@docker system prune -f
	@echo "Log files cleaned!"
