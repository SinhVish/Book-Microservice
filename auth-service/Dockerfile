FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy auth-service files
COPY auth-service/go.mod auth-service/go.sum ./auth-service/

WORKDIR /app/auth-service

# Download dependencies
RUN go mod download

# Copy auth-service source code
COPY auth-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/auth-service/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 