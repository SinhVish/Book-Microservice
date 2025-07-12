# Development Dockerfile with Air hot reload
FROM golang:1.24

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy shared module first
COPY shared/ ./shared/

# Copy auth-service files
COPY auth-service/go.mod auth-service/go.sum ./auth-service/

WORKDIR /app/auth-service

# Download dependencies
RUN go mod download

# Copy auth-service source code
COPY auth-service/ ./

# Create tmp directory for Air
RUN mkdir -p tmp

# Expose port
EXPOSE 8080

# Run with Air for hot reload
CMD ["air", "-c", ".air.toml"] 