# Root Makefile - delegates to backend
.PHONY: help dev run build test test-coverage clean health logs logs-app logs-db logs-redis docker-up docker-down docker-build docker-clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development commands (delegate to backend)
dev: ## Start development server with hot reload
	@cd backend && $(MAKE) dev

run: ## Start production server
	@cd backend && $(MAKE) run

build: ## Build the application
	@cd backend && $(MAKE) build

test: ## Run tests
	@cd backend && $(MAKE) test

test-coverage: ## Run tests with coverage
	@cd backend && $(MAKE) test-coverage

clean: ## Clean build artifacts
	@cd backend && $(MAKE) clean

# Health and monitoring
health: ## Check application health
	@cd backend && $(MAKE) health

logs: ## Show all container logs
	@cd backend && $(MAKE) logs

logs-app: ## Show application logs
	@cd backend && $(MAKE) logs-app

logs-db: ## Show database logs
	@cd backend && $(MAKE) logs-db

logs-redis: ## Show Redis logs
	@cd backend && $(MAKE) logs-redis

# Docker commands (run from root since docker-compose files are here)
docker-up: ## Start all services with Docker
	docker-compose up -d

docker-down: ## Stop all Docker services
	docker-compose down

docker-build: ## Build Docker images
	docker-compose build

docker-clean: ## Clean Docker resources
	docker-compose down -v --remove-orphans
	docker system prune -f

# Frontend commands
web-dev: ## Start frontend development server
	@cd web && npm run dev

web-build: ## Build frontend for production
	@cd web && npm run build

web-install: ## Install frontend dependencies
	@cd web && npm install

# Full stack development
full-dev: ## Start both backend and frontend in development mode
	@echo "Starting backend..."
	@cd backend && $(MAKE) dev &
	@echo "Starting frontend..."
	@cd web && npm run dev &
	@echo "Both services started. Use Ctrl+C to stop."

# Migration commands (delegate to backend)
migrate-up: ## Run database migrations
	@cd backend && $(MAKE) migrate-up

migrate-down: ## Rollback database migrations
	@cd backend && $(MAKE) migrate-down

migrate-status: ## Show migration status
	@cd backend && $(MAKE) migrate-status
