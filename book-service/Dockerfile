# Development Dockerfile with Air hot reload
FROM golang:1.24

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy shared module first
COPY shared/ ./shared/

# Copy book-service files
COPY book-service/go.mod book-service/go.sum ./book-service/

WORKDIR /app/book-service

# Download dependencies
RUN go mod download

# Copy book-service source code
COPY book-service/ ./

# Create tmp directory for Air
RUN mkdir -p tmp

# Expose port
EXPOSE 8082

# Run with Air for hot reload
CMD ["air", "-c", ".air.toml"] 