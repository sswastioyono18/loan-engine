# Testing the Loan Engine

This document provides comprehensive instructions on how to test the loan engine RESTful API.

## Prerequisites

Before testing, ensure you have:

- Docker and Docker Compose installed
- Go 1.24+ installed (for running tests locally)
- A tool to make HTTP requests (curl, Postman, Insomnia, etc.)

## Database Migrations

The loan engine includes a database migration system using Goose to manage schema changes. Migrations need to be run manually before starting the application.


## Running the Application

### 1. Using Docker Compose (Recommended)

```bash
# Start the application and database
docker-compose up -d

# Check the logs to ensure everything is running
docker-compose logs -f

# The API will be available at http://localhost:8080
# Adminer will be available at http://localhost:8081
```

### 2. Running Locally

```bash
# Install dependencies
go mod tidy

# Create a .env file with database configuration
cp .env.example .env
# Edit .env with your database credentials

# Run the application (migrations are handled separately)
go run main.go

# The API will be available at http://localhost:8080
```

## Running Migrations

After DB up, you need to run migrations, use the migration tool:

```bash
# Build the migration tool
go build -o migrate ./cmd/migrate

# Apply all pending migrations
./migrate -action up

# Check migration status
./migrate -action status

# Rollback the last migration (not recommended in production)
./migrate -action down

# Run migrations from a custom directory
./migrate -action up -dir ./my-migrations
```

## API Testing Guide

### 1. Health Check

```bash
curl -X GET http://localhost:8080/health
```

Expected response: `OK`

### 2. Complete Loan Lifecycle Test

#### Step 1: Create a Borrower

```bash
curl -X POST http://localhost:8080/api/v1/borrowers \
  -H "Content-Type: application/json" \
  -d '{
    "id_number": "B001",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "phone": "+621234567890",
    "address": "Jalan Tedeng Aling Aling"
  }'
```

#### Step 2: Create a Loan (Initial State: Proposed)

```bash
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": 1,
    "principal_amount": 1000000.00,
    "rate": 0.05,
    "roi": 0.08,
    "agreement_letter_link": "https://example.com/agreement.pdf"
  }'
```

Expected response: Loan object with `current_state: "proposed"`

#### Step 3: Approve the Loan (State: Proposed → Approved)

```bash
curl -X POST http://localhost:8080/api/v1/loans/1/approve \
  -H "Content-Type: application/json" \
  -d '{
    "field_validator_employee_id": "emp001",
    "proof_image_url": "https://example.com/proof.jpg"
  }'
```

Expected response: Success message

#### Step 4: Create an Investor

```bash
curl -X POST http://localhost:8080/api/v1/investors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Smith",
    "email": "jane.smith@example.com",
    "phone": "+0987654321",
    "investor_id": "INV001"
  }'
```

#### Step 5: Invest in the Loan (State: Approved → Invested when fully funded)

```bash
curl -X POST http://localhost:8080/api/v1/loans/1/invest \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": 1,
    "investment_amount": 1000000.00
  }'
```

Expected response: Success message

#### Step 6: Disburse the Loan (State: Invested → Disbursed)

```bash
curl -X POST http://localhost:8080/api/v1/loans/1/disburse \
  -H "Content-Type: application/json" \
  -d '{
    "field_officer_employee_id": "emp002",
    "agreement_letter_signed_url": "https://example.com/signed-agreement.pdf"
  }'
```

Expected response: Success message

### 3. Query Endpoints

#### Get Loan by ID

```bash
curl -X GET http://localhost:8080/api/v1/loans/1
```

#### List Loans

```bash
curl -X GET http://localhost:8080/api/v1/loans
```

#### List Loans by State

```bash
curl -X GET http://localhost:8080/api/v1/loans/state/approved
```

#### Get Loans with Pagination

```bash
curl -X GET "http://localhost:8080/api/v1/loans?state=proposed&offset=0&limit=10"
```

### 4. Test State Transition Validation

#### Attempt to Approve an Already Approved Loan

```bash
curl -X POST http://localhost:8080/api/v1/loans/1/approve \
  -H "Content-Type: application/json" \
  -d '{
    "field_validator_employee_id": "emp001",
    "proof_image_url": "https://example.com/proof.jpg"
  }'
```

Expected response: Error message indicating the loan is not in proposed state

#### Attempt to Invest in a Proposed Loan (should fail)

Create a new loan and try to invest before approval:

```bash
# Create another loan
curl -X POST http://localhost:8080/api/v1/loans \
  -H "Content-Type: application/json" \
  -d '{
    "borrower_id": 1,
    "principal_amount": 5000.0,
    "rate": 0.05,
    "roi": 0.08,
    "agreement_letter_link": "https://example.com/agreement2.pdf"
  }'

# Try to invest in the proposed loan (should fail)
curl -X POST http://localhost:8080/api/v1/loans/2/invest \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": 1,
    "investment_amount": 5000.0
  }'
```

Expected response: Error message indicating the loan is not in approved state

#### Attempt to Invest More Than Principal Amount

```bash
curl -X POST http://localhost:8080/api/v1/loans/1/invest \
  -H "Content-Type: application/json" \
  -d '{
    "investor_id": 1,
    "investment_amount": 15000.0
  }'
```

Expected response: Error message indicating investment exceeds remaining principal

### 5. Unit Tests

Run the unit tests to verify the service layer logic:

```bash
# Run all tests
go test ./... -v

# Run specific service tests
go test ./internal/services/... -v

# Run with coverage
go test ./... -v -cover
```

### 6. Database Schema Verification

After running the application, you can connect to the PostgreSQL database to verify the data:

```bash
# Connect to the database container
docker exec -it loan-engine-db psql -U postgres -d loan_engine

# Check migration status
SELECT * FROM schema_migrations;

# Check loan states
SELECT id, current_state, total_invested_amount FROM loans;

# Check loan approvals
SELECT * FROM loan_approvals;

# Check loan investments
SELECT * FROM loan_investments;
```

## Expected Behavior

1. **State Transitions**: Loans can only move forward in state (proposed → approved → invested → disbursed)
2. **Investment Validation**: Total investment cannot exceed the principal amount
3. **Business Logic**: Proper validation at each state transition
4. **Audit Trail**: All state changes are recorded in the loan_state_history table
5. **Email Notifications**: Mock email service logs when notifications are sent

## Troubleshooting

### Common Issues

1. **Database Connection**: Ensure PostgreSQL is running and credentials are correct
2. **Port Conflicts**: Make sure port 8080 is available
3. **Environment Variables**: Verify .env file has correct database configuration
4. **Migration Errors**: Check that migration files are properly formatted and have correct permissions

### Logs

Check application logs for errors:

```bash
# Docker logs
docker-compose logs

# For specific service
docker-compose logs api
```

## Cleanup

To stop and remove containers:

```bash
docker-compose down
```

To remove volumes as well:

```bash
docker-compose down -v