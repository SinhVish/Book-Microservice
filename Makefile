# Makefile for Go Microservices

.PHONY: help up-build logs up down clean

# Default target
help:
	@echo "Available commands:"
	@echo "  up             - Start all services with Docker Compose"
	@echo "  up-build       - Start development with hot reload (recommended)"
	@echo "  up-clean       - Total clean restart with Docker"
	@echo "  logs           - Logs from all services"
	@echo "  down           - Stop all services"
	@echo "  clean          - Clean up containers and volumes"
	@echo "  tidy           - Tidy up dependencies"

# Development with hot reload in Docker (recommended)
up-build:
	@echo "Starting development with hot reload in Docker..."
	@docker compose up --build -d
	@echo "Services started with hot reload!"
	@echo "Auth service: http://localhost:8080"
	@echo "User service: http://localhost:8081"
	@echo "Use 'make logs' to view logs"
	@echo "Use 'make down' to stop services"

# Full Docker development with hot reload
up-clean:
	@echo "Performing total clean restart..."
	@echo "Stopping all services..."
	docker compose down -v
	@echo "Pruning Docker system..."
	docker system prune -f
	@echo "Rebuilding all services..."
	docker compose build --no-cache
	@echo "Starting all services..."
	docker compose up -d
	@echo "Total clean restart complete!"

# Logs from all services
logs:
	@echo "Logs from all services (press Ctrl+C to stop)..."
	@docker compose logs -f --timestamps

# Start all services with Docker Compose
up:
	@echo "Starting all services..."
	@docker compose up -d
	@echo "Services started!"

# Stop all services
down:
	@echo "Stopping all services..."
	@docker compose down
	@echo "Services stopped!"

# Clean up containers and volumes
clean:
	@echo "Cleaning up..."
	@docker compose down -v
	@docker system prune -f
	@echo "Cleanup complete!" 

# Tidy up dependencies
tidy:
	@echo "Tidying up dependencies..."
	@cd auth-service && go mod tidy
	@cd user-service && go mod tidy
	@cd book-service && go mod tidy
	@cd shared && go mod tidy
	@echo "Dependencies tidied up!"