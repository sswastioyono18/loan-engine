# Repository Pattern Design for Loan Engine

## Overview
This document outlines the repository pattern implementation for the loan engine system using Go and SQLX. The repository pattern provides an abstraction layer between the business logic and data access layers.

## Project Structure
```
loan-engine/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── models/
│   │   ├── borrower.go
│   │   ├── loan.go
│   │   ├── loan_approval.go
│   │   ├── loan_disbursement.go
│   │   ├── investor.go
│   │   ├── loan_investment.go
│   │   ├── loan_state_history.go
│   │   └── user.go
│   ├── repositories/
│   │   ├── borrower_repository.go
│   │   ├── loan_repository.go
│   │   ├── loan_approval_repository.go
│   │   ├── loan_disbursement_repository.go
│   │   ├── investor_repository.go
│   │   ├── loan_investment_repository.go
│   │   ├── loan_state_history_repository.go
│   │   └── user_repository.go
│   ├── services/
│   │   ├── borrower_service.go
│   │   ├── loan_service.go
│   │   ├── loan_approval_service.go
│   │   ├── loan_disbursement_service.go
│   │   ├── investor_service.go
│   │   ├── loan_investment_service.go
│   │   └── auth_service.go
│   ├── handlers/
│   │   ├── borrower_handler.go
│   │   ├── loan_handler.go
│   │   ├── loan_approval_handler.go
│   │   ├── loan_disbursement_handler.go
│   │   ├── investor_handler.go
│   │   ├── loan_investment_handler.go
│   │   └── auth_handler.go
│   └── database/
│       └── db.go
├── pkg/
│   ├── middleware/
│   │   └── auth.go
│   └── utils/
│       └── helpers.go
└── go.mod
```

## Repository Interface Definitions

### 1. Borrower Repository

```go
// internal/repositories/borrower_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type BorrowerRepository interface {
    Create(ctx context.Context, borrower *models.Borrower) error
    GetByID(ctx context.Context, id int) (*models.Borrower, error)
    GetByBorrowerIDNumber(ctx context.Context, borrowerIDNumber string) (*models.Borrower, error)
    Update(ctx context.Context, borrower *models.Borrower) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, offset, limit int) ([]*models.Borrower, error)
}
```

### 2. Loan Repository

```go
// internal/repositories/loan_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type LoanRepository interface {
    Create(ctx context.Context, loan *models.Loan) error
    GetByID(ctx context.Context, id int) (*models.Loan, error)
    GetByLoanID(ctx context.Context, loanID string) (*models.Loan, error)
    Update(ctx context.Context, loan *models.Loan) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, state *string, offset, limit int) ([]*models.Loan, error)
    UpdateState(ctx context.Context, id int, newState string) error
    UpdateTotalInvestedAmount(ctx context.Context, loanID int, amount float64) error
    GetByState(ctx context.Context, state string) ([]*models.Loan, error)
}
```

### 3. Loan Approval Repository

```go
// internal/repositories/loan_approval_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type LoanApprovalRepository interface {
    Create(ctx context.Context, approval *models.LoanApproval) error
    GetByLoanID(ctx context.Context, loanID int) (*models.LoanApproval, error)
    GetByID(ctx context.Context, id int) (*models.LoanApproval, error)
    Update(ctx context.Context, approval *models.LoanApproval) error
    Delete(ctx context.Context, id int) error
}
```

### 4. Loan Disbursement Repository

```go
// internal/repositories/loan_disbursement_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type LoanDisbursementRepository interface {
    Create(ctx context.Context, disbursement *models.LoanDisbursement) error
    GetByLoanID(ctx context.Context, loanID int) (*models.LoanDisbursement, error)
    GetByID(ctx context.Context, id int) (*models.LoanDisbursement, error)
    Update(ctx context.Context, disbursement *models.LoanDisbursement) error
    Delete(ctx context.Context, id int) error
}
```

### 5. Investor Repository

```go
// internal/repositories/investor_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type InvestorRepository interface {
    Create(ctx context.Context, investor *models.Investor) error
    GetByID(ctx context.Context, id int) (*models.Investor, error)
    GetByInvestorID(ctx context.Context, investorID string) (*models.Investor, error)
    GetByEmail(ctx context.Context, email string) (*models.Investor, error)
    Update(ctx context.Context, investor *models.Investor) error
    Delete(ctx context.Context, id int) error
    List(ctx context.Context, offset, limit int) ([]*models.Investor, error)
}
```

### 6. Loan Investment Repository

```go
// internal/repositories/loan_investment_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type LoanInvestmentRepository interface {
    Create(ctx context.Context, investment *models.LoanInvestment) error
    GetByID(ctx context.Context, id int) (*models.LoanInvestment, error)
    GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanInvestment, error)
    GetByInvestorID(ctx context.Context, investorID int) ([]*models.LoanInvestment, error)
    GetByLoanAndInvestor(ctx context.Context, loanID, investorID int) (*models.LoanInvestment, error)
    Update(ctx context.Context, investment *models.LoanInvestment) error
    Delete(ctx context.Context, id int) error
    GetTotalInvestedAmountByLoan(ctx context.Context, loanID int) (float64, error)
}
```

### 7. Loan State History Repository

```go
// internal/repositories/loan_state_history_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type LoanStateHistoryRepository interface {
    Create(ctx context.Context, history *models.LoanStateHistory) error
    GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanStateHistory, error)
    GetLatestByLoanID(ctx context.Context, loanID int) (*models.LoanStateHistory, error)
    List(ctx context.Context, loanID int, offset, limit int) ([]*models.LoanStateHistory, error)
}
```

### 8. User Repository

```go
// internal/repositories/user_repository.go
package repositories

import (
    "context"
    "loan-engine/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id int) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    GetByUserID(ctx context.Context, userID string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id int) error
    UpdatePassword(ctx context.Context, id int, hashedPassword string) error
}
```

## Repository Implementation Structure

### Base Repository with Common Methods

```go
// internal/repositories/base_repository.go
package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "loan-engine/internal/database"
)

type BaseRepository struct {
    DB *sql.DB
}

func NewBaseRepository(db *database.DB) *BaseRepository {
    return &BaseRepository{
        DB: db.SqlxDB.DB,
    }
}

// Common methods for all repositories
func (r *BaseRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
    return r.DB.BeginTx(ctx, nil)
}

func (r *BaseRepository) Commit(tx *sql.Tx) error {
    return tx.Commit()
}

func (r *BaseRepository) Rollback(tx *sql.Tx) error {
    return tx.Rollback()
}
```

### Example Repository Implementation (Loan Repository)

```go
// internal/repositories/loan_repository_impl.go
package repositories

import (
    "context"
    "database/sql"
    "fmt"
    "loan-engine/internal/models"
    "github.com/jmoiron/sqlx"
)

type loanRepositoryImpl struct {
    db *sqlx.DB
}

func NewLoanRepository(db *sqlx.DB) LoanRepository {
    return &loanRepositoryImpl{db: db}
}

func (r *loanRepositoryImpl) Create(ctx context.Context, loan *models.Loan) error {
    query := `
        INSERT INTO loans (
            loan_id, borrower_id, principal_amount, rate, roi, 
            agreement_letter_link, current_state, total_invested_amount
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `
    
    err := r.db.QueryRowContext(
        ctx, query,
        loan.LoanID, loan.BorrowerID, loan.PrincipalAmount,
        loan.Rate, loan.ROI, loan.AgreementLetterLink,
        loan.CurrentState, loan.TotalInvestedAmount,
    ).Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)
    
    return err
}

func (r *loanRepositoryImpl) GetByID(ctx context.Context, id int) (*models.Loan, error) {
    query := `
        SELECT id, loan_id, borrower_id, principal_amount, rate, roi,
               agreement_letter_link, current_state, total_invested_amount,
               created_at, updated_at
        FROM loans WHERE id = $1
    `
    
    var loan models.Loan
    err := r.db.GetContext(ctx, &loan, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("loan not found")
        }
        return nil, err
    }
    
    return &loan, nil
}

func (r *loanRepositoryImpl) GetByLoanID(ctx context.Context, loanID string) (*models.Loan, error) {
    query := `
        SELECT id, loan_id, borrower_id, principal_amount, rate, roi,
               agreement_letter_link, current_state, total_invested_amount,
               created_at, updated_at
        FROM loans WHERE loan_id = $1
    `
    
    var loan models.Loan
    err := r.db.GetContext(ctx, &loan, query, loanID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("loan not found")
        }
        return nil, err
    }
    
    return &loan, nil
}

func (r *loanRepositoryImpl) Update(ctx context.Context, loan *models.Loan) error {
    query := `
        UPDATE loans SET 
            borrower_id = $1, principal_amount = $2, rate = $3, roi = $4,
            agreement_letter_link = $5, current_state = $6, 
            total_invested_amount = $7, updated_at = NOW()
        WHERE id = $8
    `
    
    result, err := r.db.ExecContext(
        ctx, query,
        loan.BorrowerID, loan.PrincipalAmount, loan.Rate, loan.ROI,
        loan.AgreementLetterLink, loan.CurrentState, loan.TotalInvestedAmount,
        loan.ID,
    )
    
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("loan not found")
    }
    
    return nil
}

func (r *loanRepositoryImpl) Delete(ctx context.Context, id int) error {
    query := "DELETE FROM loans WHERE id = $1"
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("loan not found")
    }
    
    return nil
}

func (r *loanRepositoryImpl) List(ctx context.Context, state *string, offset, limit int) ([]*models.Loan, error) {
    query := "SELECT id, loan_id, borrower_id, principal_amount, rate, roi, agreement_letter_link, current_state, total_invested_amount, created_at, updated_at FROM loans"
    args := []interface{}{}
    
    if state != nil {
        query += " WHERE current_state = $1"
        args = append(args, *state)
    }
    
    query += " ORDER BY created_at DESC LIMIT $2 OFFSET $3"
    args = append(args, limit, offset)
    
    var loans []*models.Loan
    err := r.db.SelectContext(ctx, &loans, query, args...)
    if err != nil {
        return nil, err
    }
    
    return loans, nil
}

func (r *loanRepositoryImpl) UpdateState(ctx context.Context, id int, newState string) error {
    query := "UPDATE loans SET current_state = $1, updated_at = NOW() WHERE id = $2"
    result, err := r.db.ExecContext(ctx, query, newState, id)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("loan not found")
    }
    
    return nil
}

func (r *loanRepositoryImpl) UpdateTotalInvestedAmount(ctx context.Context, loanID int, amount float64) error {
    query := "UPDATE loans SET total_invested_amount = $1, updated_at = NOW() WHERE id = $2"
    result, err := r.db.ExecContext(ctx, query, amount, loanID)
    if err != nil {
        return err
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("loan not found")
    }
    
    return nil
}

func (r *loanRepositoryImpl) GetByState(ctx context.Context, state string) ([]*models.Loan, error) {
    query := "SELECT id, loan_id, borrower_id, principal_amount, rate, roi, agreement_letter_link, current_state, total_invested_amount, created_at, updated_at FROM loans WHERE current_state = $1 ORDER BY created_at DESC"
    
    var loans []*models.Loan
    err := r.db.SelectContext(ctx, &loans, query, state)
    if err != nil {
        return nil, err
    }
    
    return loans, nil
}
```

## Repository Factory Pattern

```go
// internal/repositories/factory.go
package repositories

import (
    "loan-engine/internal/database"
    "github.com/jmoiron/sqlx"
)

type RepositoryFactory struct {
    db *sqlx.DB
}

func NewRepositoryFactory(database *database.DB) *RepositoryFactory {
    return &RepositoryFactory{
        db: database.SqlxDB,
    }
}

func (f *RepositoryFactory) BorrowerRepository() BorrowerRepository {
    return NewBorrowerRepository(f.db)
}

func (f *RepositoryFactory) LoanRepository() LoanRepository {
    return NewLoanRepository(f.db)
}

func (f *RepositoryFactory) LoanApprovalRepository() LoanApprovalRepository {
    return NewLoanApprovalRepository(f.db)
}

func (f *RepositoryFactory) LoanDisbursementRepository() LoanDisbursementRepository {
    return NewLoanDisbursementRepository(f.db)
}

func (f *RepositoryFactory) InvestorRepository() InvestorRepository {
    return NewInvestorRepository(f.db)
}

func (f *RepositoryFactory) LoanInvestmentRepository() LoanInvestmentRepository {
    return NewLoanInvestmentRepository(f.db)
}

func (f *RepositoryFactory) LoanStateHistoryRepository() LoanStateHistoryRepository {
    return NewLoanStateHistoryRepository(f.db)
}

func (f *RepositoryFactory) UserRepository() UserRepository {
    return NewUserRepository(f.db)
}
```

## Database Connection Setup

```go
// internal/database/db.go
package database

import (
    "fmt"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // PostgreSQL driver
)

type DB struct {
    SqlxDB *sqlx.DB
}

func NewDBConnection(connectionString string) (*DB, error) {
    db, err := sqlx.Connect("postgres", connectionString)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Test the connection
    if err = db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &DB{SqlxDB: db}, nil
}

func (d *DB) Close() error {
    return d.SqlxDB.Close()
}