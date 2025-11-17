# State Transition Validation in Loan Engine

## Overview
The loan engine implements strict state transition validation to ensure loans can only move forward through their lifecycle (proposed → approved → invested → disbursed) and never regress to previous states. This validation occurs at multiple levels: database constraints, application logic, and business rules.

## State Transition Rules

### Valid Transitions
- **proposed** → **approved** (when loan is approved by staff)
- **approved** → **invested** (when total investment equals principal amount)
- **invested** → **disbursed** (when loan is given to borrower)

### Invalid Transitions
- **approved** → **proposed** (cannot regress)
- **invested** → **proposed** (cannot regress)
- **invested** → **approved** (cannot regress)
- **disbursed** → **proposed/approved/invested** (cannot regress)

## Database-Level Validation

### Trigger-Based Validation
The database implements a trigger that validates state transitions at the database level:

```sql
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

CREATE TRIGGER trigger_validate_state_transition
    BEFORE INSERT ON loan_state_history
    FOR EACH ROW
    EXECUTE FUNCTION validate_state_transition();
```

This trigger ensures that even if application-level validation is bypassed, the database will prevent invalid state transitions.

## Application-Level Validation

### Service Layer Validation
The service layer implements validation logic to prevent invalid transitions:

```go
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

### Specific Transition Validations

#### Approve Loan Validation
```go
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
    
    // ... rest of the approval logic
}
```

#### Invest in Loan Validation
```go
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
    
    // ... rest of the investment logic
}
```

#### Disburse Loan Validation
```go
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
    
    // ... rest of the disbursement logic
}
```

## Audit Trail

### State History Tracking
Every state transition is recorded in the `loan_state_history` table:

```sql
CREATE TABLE loan_state_history (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id) ON DELETE CASCADE NOT NULL,
    previous_state VARCHAR(20) CHECK (previous_state IN ('proposed', 'approved', 'invested', 'disbursed')),
    new_state VARCHAR(20) NOT NULL CHECK (new_state IN ('proposed', 'approved', 'invested', 'disbursed')),
    transition_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### State Transition Recording
```go
// After successful state transition
stateHistory := &models.LoanStateHistory{
    LoanID:          loanID,
    PreviousState:   loan.CurrentState,
    NewState:        newState,
    TransitionReason: reason,
}

err = s.repoFactory.LoanStateHistoryRepository().Create(ctx, stateHistory)
if err != nil {
    tx.Rollback()
    return fmt.Errorf("failed to create state history: %w", err)
}
```

## Transaction Safety

### ACID Compliance
State transitions are wrapped in database transactions to ensure data consistency:

```go
// Start transaction
tx, err := s.repoFactory.(*repositories.RepositoryFactory).DB().BeginTx(ctx, nil)
if err != nil {
    return fmt.Errorf("failed to start transaction: %w", err)
}

// Perform state transition operations
// Update loan state
// Create state history record
// Update related entities

// Commit transaction
err = tx.Commit()
if err != nil {
    return fmt.Errorf("failed to commit transaction: %w", err)
}
```

## Error Handling

### Validation Error Types
The system provides specific error messages for different validation failures:

- `INVALID_STATE_TRANSITION`: Attempted to perform invalid state transition
- `LOAN_ALREADY_APPROVED`: Loan is already in approved state
- `LOAN_ALREADY_INVESTED`: Loan is already in invested state
- `LOAN_ALREADY_DISBURSED`: Loan is already in disbursed state

## Summary

The state transition validation in the loan engine provides multiple layers of protection:

1. **Database-level constraints** prevent invalid transitions at the persistence layer
2. **Application-level validation** provides business logic enforcement
3. **Audit trail** maintains a complete history of all state changes
4. **Transaction safety** ensures data consistency during transitions
5. **Comprehensive error handling** provides clear feedback for validation failures

This multi-layered approach ensures that the loan engine maintains data integrity and follows the business rules for loan lifecycle management while providing a robust audit trail for compliance purposes.