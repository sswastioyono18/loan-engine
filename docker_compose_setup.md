# Docker Compose Setup for Loan Engine

## Overview
This document outlines the Docker Compose configuration for the loan engine system, including PostgreSQL database, the Go application, and any necessary supporting services.

## Docker Compose File

```yaml
# docker-compose.yml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    container_name: loan_engine_postgres
    environment:
      POSTGRES_DB: loan_engine_db
      POSTGRES_USER: loan_engine_user
      POSTGRES_PASSWORD: loan_engine_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - loan_engine_network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U loan_engine_user -d loan_engine_db"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis for caching and session management
  redis:
    image: redis:7-alpine
    container_name: loan_engine_redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    networks:
      - loan_engine_network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Loan Engine Application
  loan-engine:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: loan_engine_app
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=loan_engine_user
      - DB_PASSWORD=loan_engine_password
      - DB_NAME=loan_engine_db
      - DB_SSL_MODE=disable
      - REDIS_URL=redis:6379
      - JWT_SECRET=your_jwt_secret_key_here
      - PORT=8080
      - ENV=development
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - loan_engine_network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # Adminer - Database management tool
  adminer:
    image: adminer
    container_name: loan_engine_adminer
    ports:
      - "8081:8080"
    environment:
      - ADMINER_DEFAULT_SERVER=postgres
      - ADMINER_PLUGINS=tables-filter tinymce
    depends_on:
      - postgres
    networks:
      - loan_engine_network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  loan_engine_network:
    driver: bridge
```

## Database Initialization Script

```sql
-- init.sql
-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create borrowers table
CREATE TABLE IF NOT EXISTS borrowers (
    id SERIAL PRIMARY KEY,
    borrower_id_number VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loans table
CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    loan_id VARCHAR(50) UNIQUE NOT NULL,
    borrower_id INTEGER REFERENCES borrowers(id) NOT NULL,
    principal_amount DECIMAL(15, 2) NOT NULL,
    rate DECIMAL(5, 2) NOT NULL,
    roi DECIMAL(5, 2) NOT NULL,
    agreement_letter_link TEXT,
    current_state VARCHAR(20) DEFAULT 'proposed' CHECK (current_state IN ('proposed', 'approved', 'invested', 'disbursed')),
    total_invested_amount DECIMAL(15, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loan_approvals table
CREATE TABLE IF NOT EXISTS loan_approvals (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id) ON DELETE CASCADE NOT NULL,
    field_validator_employee_id VARCHAR(50) NOT NULL,
    approval_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    proof_image_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loan_disbursements table
CREATE TABLE IF NOT EXISTS loan_disbursements (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id) ON DELETE CASCADE NOT NULL,
    field_officer_employee_id VARCHAR(50) NOT NULL,
    disbursement_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    agreement_letter_signed_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create investors table
CREATE TABLE IF NOT EXISTS investors (
    id SERIAL PRIMARY KEY,
    investor_id VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loan_investments table
CREATE TABLE IF NOT EXISTS loan_investments (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id) ON DELETE CASCADE NOT NULL,
    investor_id INTEGER REFERENCES investors(id) NOT NULL,
    investment_amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure investment amount is positive
    CONSTRAINT positive_investment CHECK (investment_amount > 0),
    
    -- Ensure no duplicate investments by same investor in same loan
    UNIQUE(loan_id, investor_id)
);

-- Create loan_state_history table
CREATE TABLE IF NOT EXISTS loan_state_history (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id) ON DELETE CASCADE NOT NULL,
    previous_state VARCHAR(20) CHECK (previous_state IN ('proposed', 'approved', 'invested', 'disbursed')),
    new_state VARCHAR(20) NOT NULL CHECK (new_state IN ('proposed', 'approved', 'invested', 'disbursed')),
    transition_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    user_type VARCHAR(20) NOT NULL CHECK (user_type IN ('staff', 'investor', 'admin')),
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX IF NOT EXISTS idx_loans_current_state ON loans(current_state);
CREATE INDEX IF NOT EXISTS idx_loans_loan_id ON loans(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_approvals_loan_id ON loan_approvals(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_disbursements_loan_id ON loan_disbursements(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_investments_loan_id ON loan_investments(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_investments_investor_id ON loan_investments(investor_id);
CREATE INDEX IF NOT EXISTS idx_loan_state_history_loan_id ON loan_state_history(loan_id);

-- Create functions and triggers
CREATE OR REPLACE FUNCTION update_loan_total_invested()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE loans 
    SET total_invested_amount = (
        SELECT COALESCE(SUM(investment_amount), 0) 
        FROM loan_investments 
        WHERE loan_id = NEW.loan_id
    )
    WHERE id = NEW.loan_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER IF NOT EXISTS trigger_update_loan_total_invested
    AFTER INSERT OR UPDATE ON loan_investments
    FOR EACH ROW
    EXECUTE FUNCTION update_loan_total_invested();

CREATE OR REPLACE FUNCTION validate_state_transition()
RETURNS TRIGGER AS $$
DECLARE
    current_state VARCHAR(20);
BEGIN
    -- Get current state of the loan
    SELECT current_state INTO current_state FROM loans WHERE id = NEW.loan_id;
    
    -- Validate state transition (only forward transitions allowed)
    IF (current_state = 'approved' AND NEW.new_state = 'proposed') OR
       (current_state = 'invested' AND (NEW.new_state = 'proposed' OR NEW.new_state = 'approved')) OR
       (current_state = 'disbursed' AND NEW.new_state IN ('proposed', 'approved', 'invested')) THEN
        RAISE EXCEPTION 'Invalid state transition from % to %', current_state, NEW.new_state;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER IF NOT EXISTS trigger_validate_state_transition
    BEFORE INSERT ON loan_state_history
    FOR EACH ROW
    EXECUTE FUNCTION validate_state_transition();
```

## Application Dockerfile

```dockerfile
# Dockerfile
# Use the official Golang image as the base image
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Change ownership of the binary to the non-root user
RUN chown appuser:appuser main

# Switch to the non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
```

## Environment Configuration

```bash
# .env
# Database Configuration
DB_HOST=postgres
DB_PORT=5432
DB_USER=loan_engine_user
DB_PASSWORD=loan_engine_password
DB_NAME=loan_engine_db
DB_SSL_MODE=disable

# Application Configuration
PORT=8080
ENV=development
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRY_HOURS=24

# Redis Configuration
REDIS_URL=redis:6379

# External Services (for mocking in development)
EMAIL_SERVICE_URL=http://mock-email-service:3000
STORAGE_SERVICE_URL=http://mock-storage-service:3000

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
```

## Development Docker Compose Override

```yaml
# docker-compose.dev.yml
version: '3.8'

services:
  loan-engine:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - go-mod-cache:/go/pkg/mod
    environment:
      - ENV=development
      - LOG_LEVEL=debug
    command: [
      "sh", "-c", 
      "go mod download && go run cmd/server/main.go"
    ]

volumes:
  go-mod-cache:
```

## Dockerfile for Development

```dockerfile
# Dockerfile.dev
FROM golang:1.21-alpine

# Install dependencies for development
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Install air for hot reloading
RUN go install github.com/cosmtrek/air@latest

# Expose port
EXPOSE 8080

# Command to run with hot reloading
CMD ["air", "-c", ".air.toml"]
```

## Air Configuration for Hot Reloading

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
  bin = "tmp/main"
  cmd = "go build -o tmp/main cmd/server/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "docker"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = true
```

## Makefile for Common Operations

```makefile
# Makefile
.PHONY: help build up down logs test migrate

# Show help
help:
	@echo "Usage:"
	@echo "  make build     - Build the Docker images"
	@echo "  make up        - Start the services"
	@echo "  make down      - Stop the services"
	@echo "  make logs      - View logs"
	@echo "  make test      - Run tests"
	@echo "  make migrate   - Run database migrations"

# Build Docker images
build:
	docker-compose build

# Start services
up:
	docker-compose up -d

# Stop services
down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Run tests
test:
	docker-compose exec loan-engine go test ./...

# Run database migrations (if needed separately)
migrate:
	@echo "Database migrations run automatically on startup via init.sql"