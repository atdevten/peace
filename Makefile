.PHONY: help build run dev clean logs ps stop restart

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build the application Docker image"
	@echo "  run       - Run the full stack (production)"
	@echo "  dev       - Run development environment with hot reload"
	@echo "  clean     - Stop and remove all containers and volumes"
	@echo "  stop      - Stop all services"
	@echo "  restart   - Restart all services"
	@echo "  health    - Check service health"

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
