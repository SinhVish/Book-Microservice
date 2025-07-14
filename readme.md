# Microservice Example

## Quick Start

1. Copy environment file: `cp env.example .env`
2. Start all services: `make up-build`
3. Services will be available at:
   - Auth Service: http://localhost:8080
   - User Service: http://localhost:8081  
   - Book Service: http://localhost:8082

## Commands

- `make up-build` - Start with hot reload (recommended for development)
- `make down` - Stop all services
- `make logs` - View logs from all services
