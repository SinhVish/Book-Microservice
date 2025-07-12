# Development Dockerfile with Air hot reload
FROM golang:1.24

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy shared module first
COPY shared/ ./shared/

# Copy user-service files
COPY user-service/go.mod user-service/go.sum ./user-service/

WORKDIR /app/user-service

# Download dependencies
RUN go mod download

# Copy user-service source code
COPY user-service/ ./

# Create tmp directory for Air
RUN mkdir -p tmp

# Expose ports
EXPOSE 8081
EXPOSE 9081

# Run with Air for hot reload
CMD ["air", "-c", ".air.toml"] 