# Makefile for Go Microservices

.PHONY: help build run test clean up down logs

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build all services"
	@echo "  run            - Run services locally"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean up containers and volumes"
	@echo "  up             - Start all services with Docker Compose"
	@echo "  down           - Stop all services"
	@echo "  logs           - View logs"
	@echo "  logs-auth      - View auth service logs"
	@echo "  logs-user      - View user service logs"
	@echo "  logs-mysql     - View MySQL logs"
	@echo "  health         - Check health of all services"
	@echo "  rebuild        - Rebuild and restart services"

# Build services
build:
	@echo "Building auth service..."
	cd auth-service && go build -o ../bin/auth-service cmd/main.go
	@echo "Building user service..."
	cd user-service && go build -o ../bin/user-service cmd/main.go
	@echo "Build complete!"

# Run services locally (requires MySQL running)
run:
	@echo "Starting services locally..."
	@echo "Make sure MySQL is running: docker compose up -d mysql"
	cd auth-service && go run cmd/main.go &
	cd user-service && go run cmd/main.go &
	@echo "Services started!"

# Test services
test:
	@echo "Running tests for auth service..."
	cd auth-service && go test ./...
	@echo "Running tests for user service..."
	cd user-service && go test ./...

# Clean up
clean:
	@echo "Cleaning up..."
	docker compose down -v
	docker system prune -f
	rm -rf bin/
	@echo "Cleanup complete!"

# Start all services
up:
	@echo "Starting all services..."
	docker compose up -d
	@echo "Services started! Check health with: make health"

# Stop all services
down:
	@echo "Stopping all services..."
	docker compose down
	@echo "Services stopped!"

# View logs
logs:
	docker compose logs -f

# View auth service logs
logs-auth:
	docker compose logs -f auth-service

# View user service logs
logs-user:
	docker compose logs -f user-service

# View MySQL logs
logs-mysql:
	docker compose logs -f mysql

# Health check
health:
	@echo "Checking service health..."
	@curl -s http://localhost:8080/health || echo "Auth service not responding"
	@curl -s http://localhost:8081/health || echo "User service not responding"

# Rebuild and restart
rebuild:
	@echo "Rebuilding and restarting services..."
	docker compose down
	docker compose build --no-cache
	docker compose up -d
	@echo "Services rebuilt and restarted!"

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	mkdir -p bin
	go mod tidy -C auth-service
	go mod tidy -C user-service
	@echo "Development setup complete!"

# Quick test of the system
quick-test: up
	@echo "Waiting for services to start..."
	sleep 10
	@echo "Testing health endpoints..."
	curl -s http://localhost:8080/health | jq .
	curl -s http://localhost:8081/health | jq . 