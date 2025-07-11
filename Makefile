# Makefile for Go Microservices

.PHONY: help build run test clean up down logs prune clean-all restart status db-check test-api

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build all services"
	@echo "  run            - Run services locally"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean up containers and volumes"
	@echo "  prune          - Clean up Docker system (images, cache, etc.)"
	@echo "  clean-all      - Total clean restart with database recreation"
	@echo "  restart        - Quick restart of services"
	@echo "  up             - Start all services with Docker Compose"
	@echo "  down           - Stop all services"
	@echo "  logs           - View logs from all services"
	@echo "  logs-auth      - View auth service logs"
	@echo "  logs-user      - View user service logs"
	@echo "  logs-mysql     - View MySQL logs"
	@echo "  logs-follow    - Follow logs from all services"
	@echo "  health         - Check health of all services"
	@echo "  status         - Show running containers"
	@echo "  db-check       - Check database status and tables"
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

# View logs (last 50 lines)
logs:
	docker compose logs --tail=50

# Follow logs in real-time
logs-follow:
	docker compose logs -f

# View auth service logs
logs-auth:
	docker compose logs --tail=50 auth-service

# View user service logs  
logs-user:
	docker compose logs --tail=50 user-service

# View MySQL logs
logs-mysql:
	docker compose logs --tail=50 mysql

# Follow specific service logs
logs-auth-follow:
	docker compose logs -f auth-service

logs-user-follow:
	docker compose logs -f user-service

logs-mysql-follow:
	docker compose logs -f mysql

# Health check
health:
	@echo "Checking service health..."
	@curl -s http://localhost:8080/health || echo "Auth service not responding"
	@curl -s http://localhost:8081/health || echo "User service not responding"

# Docker system prune
prune:
	@echo "Pruning Docker system..."
	docker system prune -f
	@echo "Docker system pruned!"

# Total clean restart with database recreation
clean-all:
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

# Quick restart of services
restart:
	@echo "Restarting services..."
	docker compose restart
	@echo "Services restarted!"

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

# Show running containers
status:
	@echo "Docker containers status:"
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Check database status
db-check:
	@echo "Checking database status..."
	docker exec microservice-mysql mysql -u root -p1234 -e "SHOW DATABASES;" 2>/dev/null || echo "Database not accessible"
	@echo "Checking tables in auth_db:"
	docker exec microservice-mysql mysql -u root -p1234 -e "USE auth_db; SHOW TABLES;" 2>/dev/null || echo "auth_db not accessible"
	@echo "Checking tables in user_db:"
	docker exec microservice-mysql mysql -u root -p1234 -e "USE user_db; SHOW TABLES;" 2>/dev/null || echo "user_db not accessible"

# Quick test of the system
quick-test: up
	@echo "Waiting for services to start..."
	sleep 10
	@echo "Testing health endpoints..."
	curl -s http://localhost:8080/health | jq .
	curl -s http://localhost:8081/health | jq . 