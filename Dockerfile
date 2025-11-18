# Dockerfile
# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install git and other dependencies needed for Go modules
RUN apk add --no-cache git ca-certificates build-base

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and postgresql-client for goose
RUN apk --no-cache add ca-certificates postgresql-client

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

# Copy the migrations directory
COPY --from=builder /app/migrations/ ./migrations/

# Change ownership of the binaries to the non-root user
RUN chown appuser:appuser main migrate

# Switch to the non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]