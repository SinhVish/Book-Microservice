FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy user-service files
COPY user-service/go.mod user-service/go.sum ./user-service/

WORKDIR /app/user-service

# Download dependencies
RUN go mod download

# Copy user-service source code
COPY user-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/user-service/main .

# Expose port
EXPOSE 8081

# Run the application
CMD ["./main"] 