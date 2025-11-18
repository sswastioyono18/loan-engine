# Loan Engine

A RESTful API for managing loan applications with state transitions from proposed to approved, invested, and disbursed.

## Features

- Complete loan lifecycle management
- State transition validation
- Investor management
- Email notifications for investors
- File storage for documents and proofs
- Comprehensive audit trail

## Tech Stack

- Go (Golang)
- PostgreSQL
- Chi router
- SQLX for database operations
- Repository pattern
- Service layer architecture

## Project Structure

```
loan-engine/
├── cmd/
├── internal/
│   ├── models/           # Data models
│   ├── repositories/     # Repository interfaces and implementations
│   ├── services/         # Business logic
│   └── handlers/         # HTTP handlers
├── pkg/
│   ├── external/         # External service interfaces and mocks
│   └── util/             # Utility functions
├── docker-compose.yml
├── Dockerfile
└── go.mod
```

## API Endpoints

### Borrowers

- `POST /api/v1/borrowers` - Create a new borrower
- `GET /api/v1/borrowers/{id}` - Get borrower by ID
- `PUT /api/v1/borrowers/{id}` - Update borrower
- `DELETE /api/v1/borrowers/{id}` - Delete borrower
- `GET /api/v1/borrowers` - List borrowers

### Loans

- `POST /api/v1/loans` - Create a new loan
- `GET /api/v1/loans/{id}` - Get loan by ID
- `PUT /api/v1/loans/{id}` - Update loan
- `DELETE /api/v1/loans/{id}` - Delete loan
- `GET /api/v1/loans` - List loans with optional state filter
- `POST /api/v1/loans/{id}/approve` - Approve a loan
- `POST /api/v1/loans/{id}/invest` - Invest in a loan
- `POST /api/v1/loans/{id}/disburse` - Disburse a loan
- `GET /api/v1/loans/state/{state}` - Get loans by state

### Investors

- `POST /api/v1/investors` - Create a new investor
- `GET /api/v1/investors/{id}` - Get investor by ID
- `PUT /api/v1/investors/{id}` - Update investor
- `DELETE /api/v1/investors/{id}` - Delete investor
- `GET /api/v1/investors` - List investors

## Database Schema

The database schema includes tables for:
- Borrowers
- Loans
- Loan approvals
- Loan disbursements
- Investors
- Loan investments
- Loan state history
- Users

## State Transitions

Loans follow a strict state transition pattern:
1. `proposed` (initial state)
2. `approved` (after field validation)
3. `invested` (after sufficient investment)
4. `disbursed` (after funds are given to borrower)

State transitions can only move forward, never backward.

## Getting Started

### Prerequisites

- Go 1.19+
- Docker and Docker Compose
- PostgreSQL (if not using Docker)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/sswastioyono18/loan-engine.git
   cd loan-engine
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your database configuration
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

### Using Docker

1. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

## Configuration

The application uses environment variables for configuration:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSL_MODE`: SSL mode for database connection
- `PORT`: Application port

## Testing

Run the tests:
```bash
go test ./...
```

## API Documentation

For detailed API documentation, refer to the `api_endpoints.md` file in the project root.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License.