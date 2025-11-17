# Service Layer Design for Loan Engine

## Overview
This document outlines the service layer implementation for the loan engine system. The service layer contains the business logic and orchestrates operations between the handlers and repositories.

## Service Interface Definitions

### 1. Borrower Service

```go
// internal/services/borrower_service.go
package services

import (
    "context"
    "loan-engine/internal/models"
)

type BorrowerService interface {
    CreateBorrower(ctx context.Context, borrower *models.Borrower) error
    GetBorrowerByID(ctx context.Context, id int) (*models.Borrower, error)
    GetBorrowerByBorrowerIDNumber(ctx context.Context, borrowerIDNumber string) (*models.Borrower, error)
    UpdateBorrower(ctx context.Context, id int, borrower *models.Borrower) error
    DeleteBorrower(ctx context.Context, id int) error
    ListBorrowers(ctx context.Context, offset, limit int) ([]*models.Borrower, error)
}
```

### 2. Loan Service

```go
// internal/services/loan_service.go
package services

import (
    "context"
    "loan-engine/internal/models"
)

type LoanService interface {
    CreateLoan(ctx context.Context, loan *models.Loan) error
    GetLoanByID(ctx context.Context, id int) (*models.Loan, error)
    GetLoanByLoanID(ctx context.Context, loanID string) (*models.Loan, error)
    UpdateLoan(ctx context.Context, id int, loan *models.Loan) error
    DeleteLoan(ctx context.Context, id int) error
    ListLoans(ctx context.Context, state *string, offset, limit int) ([]*models.Loan, error)
    GetLoansByState(ctx context.Context, state string) ([]*models.Loan, error)
    
    // State transition methods
    ApproveLoan(ctx context.Context, loanID int, approvalData *models.LoanApproval) error
    InvestInLoan(ctx context.Context, loanID int, investment *models.LoanInvestment) error
    DisburseLoan(ctx context.Context, loanID int, disbursementData *models.LoanDisbursement) error
    
    // Helper methods
    GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error)
    CanTransitionToState(ctx context.Context, loanID int, newState string) (bool, error)
}
```

### 3. Loan Approval Service

```go
// internal/services/loan_approval_service.go
package services

import (
    "context"
    "loan-engine/internal/models"
)

type LoanApprovalService interface {
    CreateApproval(ctx context.Context, approval *models.LoanApproval) error
    GetApprovalByLoanID(ctx context.Context, loanID int) (*models.LoanApproval, error)
    GetApprovalByID(ctx context.Context, id int) (*models.LoanApproval, error)
}
```

### 4. Loan Disbursement Service

```go
// internal/services/loan_disbursement_service.go
package services

import (
    "context"
    "loan-engine/internal/models"
)

type LoanDisbursementService interface {
    CreateDisbursement(ctx context.Context, disbursement *models.LoanDisbursement) error
    GetDisbursementByLoanID(ctx context.Context, loanID int) (*models.LoanDisbursement, error)
    GetDisbursementByID(ctx context.Context, id int) (*models.LoanDisbursement, error)
}
```

### 5. Investor Service

```go
// internal/services/investor_service.go
package services

import (
    "context"
    "loan-engine/internal/models"
)

type InvestorService interface {
    CreateInvestor(ctx context.Context, investor *models.Investor) error
    GetInvestorByID(ctx context.Context, id int) (*models.Investor, error)
    GetInvestorByInvestorID(ctx context.Context, investorID string) (*models.Investor, error)
    GetInvestorByEmail(ctx context.Context, email string) (*models.Investor, error)
    UpdateInvestor(ctx context.Context, id int, investor *models.Investor) error
    DeleteInvestor(ctx context.Context, id int) error
    ListInvestors(ctx context.Context, offset, limit int) ([]*models.Investor, error)
}
```

### 6. Loan Investment Service

```go
// internal/services/loan_investment_service.go
package services

import (
    "context"
    "loan-engine/internal/models"
)

type LoanInvestmentService interface {
    CreateInvestment(ctx context.Context, investment *models.LoanInvestment) error
    GetInvestmentByID(ctx context.Context, id int) (*models.LoanInvestment, error)
    GetInvestmentsByLoanID(ctx context.Context, loanID int) ([]*models.LoanInvestment, error)
    GetInvestmentsByInvestorID(ctx context.Context, investorID int) ([]*models.LoanInvestment, error)
    GetInvestmentByLoanAndInvestor(ctx context.Context, loanID, investorID int) (*models.LoanInvestment, error)
}
```

### 7. Auth Service

```go
// internal/services/auth_service.go
package services

import (
    "context"
    "loan-engine/internal/models"
)

type AuthService interface {
    RegisterUser(ctx context.Context, user *models.User, password string) error
    LoginUser(ctx context.Context, email, password string) (string, error)
    RefreshToken(ctx context.Context, refreshToken string) (string, error)
    ValidateToken(ctx context.Context, token string) (*models.User, error)
    HashPassword(password string) (string, error)
    CheckPasswordHash(password, hash string) bool
}
```

## Service Implementation Structure

### Base Service with Common Methods

```go
// internal/services/base_service.go
package services

import (
    "context"
    "loan-engine/internal/repositories"
)

type BaseService struct {
    RepoFactory *repositories.RepositoryFactory
}

func NewBaseService(repoFactory *repositories.RepositoryFactory) *BaseService {
    return &BaseService{
        RepoFactory: repoFactory,
    }
}
```

### Example Service Implementation (Loan Service)

```go
// internal/services/loan_service_impl.go
package services

import (
    "context"
    "errors"
    "fmt"
    "loan-engine/internal/models"
    "loan-engine/internal/repositories"
)

type loanServiceImpl struct {
    repoFactory *repositories.RepositoryFactory
}

func NewLoanService(repoFactory *repositories.RepositoryFactory) LoanService {
    return &loanServiceImpl{
        repoFactory: repoFactory,
    }
}

func (s *loanServiceImpl) CreateLoan(ctx context.Context, loan *models.Loan) error {
    // Validate required fields
    if loan.PrincipalAmount <= 0 {
        return errors.New("principal amount must be greater than 0")
    }
    
    if loan.Rate < 0 || loan.Rate > 100 {
        return errors.New("rate must be between 0 and 100")
    }
    
    if loan.ROI < 0 || loan.ROI > 100 {
        return errors.New("ROI must be between 0 and 100")
    }
    
    // Set initial state to proposed
    loan.CurrentState = "proposed"
    loan.TotalInvestedAmount = 0.0
    
    // Create loan in repository
    return s.repoFactory.LoanRepository().Create(ctx, loan)
}

func (s *loanServiceImpl) GetLoanByID(ctx context.Context, id int) (*models.Loan, error) {
    return s.repoFactory.LoanRepository().GetByID(ctx, id)
}

func (s *loanServiceImpl) GetLoanByLoanID(ctx context.Context, loanID string) (*models.Loan, error) {
    return s.repoFactory.LoanRepository().GetByLoanID(ctx, loanID)
}

func (s *loanServiceImpl) UpdateLoan(ctx context.Context, id int, loan *models.Loan) error {
    // Get existing loan to check state
    existingLoan, err := s.repoFactory.LoanRepository().GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    // Prevent modification of certain fields based on state
    if existingLoan.CurrentState != "proposed" {
        // Only allow updating specific fields after loan is approved
        loan.BorrowerID = existingLoan.BorrowerID
        loan.PrincipalAmount = existingLoan.PrincipalAmount
        loan.Rate = existingLoan.Rate
        loan.ROI = existingLoan.ROI
        loan.AgreementLetterLink = existingLoan.AgreementLetterLink
    }
    
    return s.repoFactory.LoanRepository().Update(ctx, loan)
}

func (s *loanServiceImpl) DeleteLoan(ctx context.Context, id int) error {
    // Check if loan can be deleted (must be in proposed state)
    loan, err := s.repoFactory.LoanRepository().GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    if loan.CurrentState != "proposed" {
        return errors.New("loan can only be deleted in proposed state")
    }
    
    return s.repoFactory.LoanRepository().Delete(ctx, id)
}

func (s *loanServiceImpl) ListLoans(ctx context.Context, state *string, offset, limit int) ([]*models.Loan, error) {
    return s.repoFactory.LoanRepository().List(ctx, state, offset, limit)
}

func (s *loanServiceImpl) GetLoansByState(ctx context.Context, state string) ([]*models.Loan, error) {
    return s.repoFactory.LoanRepository().GetByState(ctx, state)
}

func (s *loanServiceImpl) ApproveLoan(ctx context.Context, loanID int, approvalData *models.LoanApproval) error {
    // Get the loan
    loan, err := s.repoFactory.LoanRepository().GetByID(ctx, loanID)
    if err != nil {
        return fmt.Errorf("loan not found: %w", err)
    }
    
    // Check if loan is in proposed state
    if loan.CurrentState != "proposed" {
        return errors.New("loan must be in proposed state to be approved")
    }
    
    // Validate approval data
    if approvalData.FieldValidatorEmployeeID == "" {
        return errors.New("field validator employee ID is required")
    }
    
    if approvalData.ProofImageUrl == "" {
        return errors.New("proof image URL is required")
    }
    
    // Start transaction
    tx, err := s.repoFactory.(*repositories.RepositoryFactory).DB().BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to start transaction: %w", err)
    }
    
    // Create loan approval record
    approvalData.LoanID = loanID
    err = s.repoFactory.LoanApprovalRepository().Create(ctx, approvalData)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create loan approval: %w", err)
    }
    
    // Update loan state to approved
    err = s.repoFactory.LoanRepository().UpdateState(ctx, loanID, "approved")
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update loan state: %w", err)
    }
    
    // Add state transition to history
    stateHistory := &models.LoanStateHistory{
        LoanID:       loanID,
        PreviousState: loan.CurrentState,
        NewState:     "approved",
        TransitionReason: "Loan approved by staff",
    }
    
    err = s.repoFactory.LoanStateHistoryRepository().Create(ctx, stateHistory)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create state history: %w", err)
    }
    
    // Commit transaction
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

func (s *loanServiceImpl) InvestInLoan(ctx context.Context, loanID int, investment *models.LoanInvestment) error {
    // Get the loan
    loan, err := s.repoFactory.LoanRepository().GetByID(ctx, loanID)
    if err != nil {
        return fmt.Errorf("loan not found: %w", err)
    }
    
    // Check if loan is in approved state
    if loan.CurrentState != "approved" {
        return errors.New("loan must be in approved state to receive investments")
    }
    
    // Validate investment amount
    if investment.InvestmentAmount <= 0 {
        return errors.New("investment amount must be greater than 0")
    }
    
    // Check if investment amount exceeds remaining principal
    remainingPrincipal := loan.PrincipalAmount - loan.TotalInvestedAmount
    if investment.InvestmentAmount > remainingPrincipal {
        return fmt.Errorf("investment amount exceeds remaining principal. Remaining: %f", remainingPrincipal)
    }
    
    // Check if investor already invested in this loan
    existingInvestment, err := s.repoFactory.LoanInvestmentRepository().GetByLoanAndInvestor(ctx, loanID, investment.InvestorID)
    if err == nil && existingInvestment != nil {
        return errors.New("investor already invested in this loan")
    }
    
    // Start transaction
    tx, err := s.repoFactory.(*repositories.RepositoryFactory).DB().BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to start transaction: %w", err)
    }
    
    // Create investment record
    investment.LoanID = loanID
    err = s.repoFactory.LoanInvestmentRepository().Create(ctx, investment)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create investment: %w", err)
    }
    
    // Update total invested amount in loan
    newTotal := loan.TotalInvestedAmount + investment.InvestmentAmount
    err = s.repoFactory.LoanRepository().UpdateTotalInvestedAmount(ctx, loanID, newTotal)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update total invested amount: %w", err)
    }
    
    // Check if loan is fully invested
    if newTotal >= loan.PrincipalAmount {
        // Update loan state to invested
        err = s.repoFactory.LoanRepository().UpdateState(ctx, loanID, "invested")
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to update loan state: %w", err)
        }
        
        // Add state transition to history
        stateHistory := &models.LoanStateHistory{
            LoanID:       loanID,
            PreviousState: loan.CurrentState,
            NewState:     "invested",
            TransitionReason: "Loan fully invested",
        }
        
        err = s.repoFactory.LoanStateHistoryRepository().Create(ctx, stateHistory)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("failed to create state history: %w", err)
        }
    }
    
    // Commit transaction
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

func (s *loanServiceImpl) DisburseLoan(ctx context.Context, loanID int, disbursementData *models.LoanDisbursement) error {
    // Get the loan
    loan, err := s.repoFactory.LoanRepository().GetByID(ctx, loanID)
    if err != nil {
        return fmt.Errorf("loan not found: %w", err)
    }
    
    // Check if loan is in invested state
    if loan.CurrentState != "invested" {
        return errors.New("loan must be in invested state to be disbursed")
    }
    
    // Check if total invested amount equals principal amount
    if loan.TotalInvestedAmount != loan.PrincipalAmount {
        return errors.New("total invested amount must equal principal amount for disbursement")
    }
    
    // Validate disbursement data
    if disbursementData.FieldOfficerEmployeeID == "" {
        return errors.New("field officer employee ID is required")
    }
    
    if disbursementData.AgreementLetterSignedUrl == "" {
        return errors.New("signed agreement letter URL is required")
    }
    
    // Start transaction
    tx, err := s.repoFactory.(*repositories.RepositoryFactory).DB().BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to start transaction: %w", err)
    }
    
    // Create loan disbursement record
    disbursementData.LoanID = loanID
    err = s.repoFactory.LoanDisbursementRepository().Create(ctx, disbursementData)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create loan disbursement: %w", err)
    }
    
    // Update loan state to disbursed
    err = s.repoFactory.LoanRepository().UpdateState(ctx, loanID, "disbursed")
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to update loan state: %w", err)
    }
    
    // Add state transition to history
    stateHistory := &models.LoanStateHistory{
        LoanID:       loanID,
        PreviousState: loan.CurrentState,
        NewState:     "disbursed",
        TransitionReason: "Loan disbursed to borrower",
    }
    
    err = s.repoFactory.LoanStateHistoryRepository().Create(ctx, stateHistory)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to create state history: %w", err)
    }
    
    // Commit transaction
    err = tx.Commit()
    if err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

func (s *loanServiceImpl) GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error) {
    return s.repoFactory.LoanInvestmentRepository().GetTotalInvestedAmountByLoan(ctx, loanID)
}

func (s *loanServiceImpl) CanTransitionToState(ctx context.Context, loanID int, newState string) (bool, error) {
    loan, err := s.repoFactory.LoanRepository().GetByID(ctx, loanID)
    if err != nil {
        return false, err
    }
    
    currentState := loan.CurrentState
    
    // Define valid state transitions
    validTransitions := map[string][]string{
        "proposed": {"approved"},
        "approved": {"invested"},
        "invested": {"disbursed"},
        "disbursed": {}, // No further transitions allowed
    }
    
    validStates, exists := validTransitions[currentState]
    if !exists {
        return false, fmt.Errorf("invalid current state: %s", currentState)
    }
    
    for _, state := range validStates {
        if state == newState {
            return true, nil
        }
    }
    
    return false, nil
}
```

## Service Factory Pattern

```go
// internal/services/factory.go
package services

import (
    "loan-engine/internal/repositories"
)

type ServiceFactory struct {
    RepoFactory *repositories.RepositoryFactory
}

func NewServiceFactory(repoFactory *repositories.RepositoryFactory) *ServiceFactory {
    return &ServiceFactory{
        RepoFactory: repoFactory,
    }
}

func (f *ServiceFactory) BorrowerService() BorrowerService {
    return NewBorrowerService(f.RepoFactory)
}

func (f *ServiceFactory) LoanService() LoanService {
    return NewLoanService(f.RepoFactory)
}

func (f *ServiceFactory) LoanApprovalService() LoanApprovalService {
    return NewLoanApprovalService(f.RepoFactory)
}

func (f *ServiceFactory) LoanDisbursementService() LoanDisbursementService {
    return NewLoanDisbursementService(f.RepoFactory)
}

func (f *ServiceFactory) InvestorService() InvestorService {
    return NewInvestorService(f.RepoFactory)
}

func (f *ServiceFactory) LoanInvestmentService() LoanInvestmentService {
    return NewLoanInvestmentService(f.RepoFactory)
}

func (f *ServiceFactory) AuthService() AuthService {
    return NewAuthService(f.RepoFactory)
}
```

## Email Service Interface (for mocking)

```go
// internal/services/email_service.go
package services

type EmailService interface {
    SendInvestmentConfirmation(investorEmail, agreementLink string) error
    SendDisbursementNotification(borrowerEmail, loanDetails string) error
}
```

## Mock Email Service Implementation

```go
// internal/services/email_service_mock.go
package services

import (
    "fmt"
    "log"
)

type MockEmailService struct{}

func NewMockEmailService() *MockEmailService {
    return &MockEmailService{}
}

func (m *MockEmailService) SendInvestmentConfirmation(investorEmail, agreementLink string) error {
    log.Printf("MOCK: Sending investment confirmation to %s with agreement link: %s", investorEmail, agreementLink)
    // In a real implementation, this would send an actual email
    return nil
}

func (m *MockEmailService) SendDisbursementNotification(borrowerEmail, loanDetails string) error {
    log.Printf("MOCK: Sending disbursement notification to %s with details: %s", borrowerEmail, loanDetails)
    // In a real implementation, this would send an actual email
    return nil
}