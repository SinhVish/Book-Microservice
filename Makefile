# Makefile for Go Microservices

.PHONY: help clean-dev clean-and-restart dev live-logs up down clean

# Default target
help:
	@echo "Available commands:"
	@echo "  dev            - Start development with hot reload (recommended)"
	@echo "  clean-dev      - Clean build with hot reload using .air.toml (local services)"
	@echo "  clean-and-restart - Total clean restart with Docker"
	@echo "  live-logs      - Live logs from all services"
	@echo "  up             - Start all services with Docker Compose"
	@echo "  down           - Stop all services"
	@echo "  clean          - Clean up containers and volumes"

# Development with hot reload in Docker (recommended)
dev:
	@echo "Starting development with hot reload in Docker..."
	@docker compose up --build -d
	@echo "Services started with hot reload!"
	@echo "Auth service: http://localhost:8080"
	@echo "User service: http://localhost:8081"
	@echo "Use 'make live-logs' to view logs"
	@echo "Use 'make down' to stop services"

# Clean build with hot reload using .air.toml
clean-dev:
	@echo "Clean build with hot reload using .air.toml..."
	@echo "Stopping and cleaning up..."
	@pkill -f "air.*auth-service" || true
	@pkill -f "air.*user-service" || true
	@docker compose down -v
	@docker system prune -f
	@echo "Installing air for hot reload..."
	@go install github.com/air-verse/air@latest
	@echo "Starting databases..."
	@docker compose up -d auth-mysql user-mysql
	@echo "Waiting for databases to be ready..."
	@sleep 5
	@echo "Starting services with hot reload..."
	@echo "Auth service: http://localhost:8080, User service: http://localhost:8081"
	@echo "Press Ctrl+C to stop all services"
	@trap 'echo "Stopping development servers..."; pkill -f "air.*auth-service"; pkill -f "air.*user-service"; docker compose down; exit' INT; \
	cd auth-service && DB_PORT=3306 DB_NAME=auth_db $$(go env GOPATH)/bin/air -c .air.toml & \
	cd user-service && DB_PORT=3307 DB_NAME=user_db $$(go env GOPATH)/bin/air -c .air.toml & \
	wait

	
# Full Docker development with hot reload
clean-and-restart:
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

# Live logs from all services
live-logs:
	@echo "Live logs from all services (press Ctrl+C to stop)..."
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