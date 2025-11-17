# Loan Engine Database Schema

## Overview
This document outlines the database schema for the loan engine system with support for loan lifecycle management, investor tracking, and state transitions.

## Database Tables

### 1. borrowers
Stores information about loan borrowers.

```sql
CREATE TABLE borrowers (
    id SERIAL PRIMARY KEY,
    borrower_id_number VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 2. loans
Main table for loan information with state tracking.

```sql
CREATE TABLE loans (
    id SERIAL PRIMARY KEY,
    loan_id VARCHAR(50) UNIQUE NOT NULL,
    borrower_id INTEGER REFERENCES borrowers(id) NOT NULL,
    principal_amount DECIMAL(15, 2) NOT NULL,
    rate DECIMAL(5, 2) NOT NULL, -- Interest rate percentage
    roi DECIMAL(5, 2) NOT NULL, -- Return of investment percentage
    agreement_letter_link TEXT,
    current_state VARCHAR(20) DEFAULT 'proposed' CHECK (current_state IN ('proposed', 'approved', 'invested', 'disbursed')),
    total_invested_amount DECIMAL(15, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 3. loan_approvals
Stores approval information when a loan transitions to approved state.

```sql
CREATE TABLE loan_approvals (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id) ON DELETE CASCADE NOT NULL,
    field_validator_employee_id VARCHAR(50) NOT NULL,
    approval_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    proof_image_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 4. loan_disbursements
Stores disbursement information when a loan transitions to disbursed state.

```sql
CREATE TABLE loan_disbursements (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER REFERENCES loans(id) ON DELETE CASCADE NOT NULL,
    field_officer_employee_id VARCHAR(50) NOT NULL,
    disbursement_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    agreement_letter_signed_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 5. investors
Stores information about investors.

```sql
CREATE TABLE investors (
    id SERIAL PRIMARY KEY,
    investor_id VARCHAR(50) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### 6. loan_investments
Tracks individual investments in loans.

```sql
CREATE TABLE loan_investments (
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
```

### 7. loan_state_history
Tracks all state transitions for audit purposes.

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

### 8. users
For authentication of system users (staff, investors).

```sql
CREATE TABLE users (
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
```

## Indexes for Performance

```sql
-- Indexes for loans table
CREATE INDEX idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX idx_loans_current_state ON loans(current_state);
CREATE INDEX idx_loans_loan_id ON loans(loan_id);

-- Indexes for loan_approvals table
CREATE INDEX idx_loan_approvals_loan_id ON loan_approvals(loan_id);

-- Indexes for loan_disbursements table
CREATE INDEX idx_loan_disbursements_loan_id ON loan_disbursements(loan_id);

-- Indexes for loan_investments table
CREATE INDEX idx_loan_investments_loan_id ON loan_investments(loan_id);
CREATE INDEX idx_loan_investments_investor_id ON loan_investments(investor_id);

-- Indexes for loan_state_history table
CREATE INDEX idx_loan_state_history_loan_id ON loan_state_history(loan_id);
```

## Constraints and Business Rules

1. **State Transition Validation**: Loans can only move forward in state (proposed → approved → invested → disbursed)
2. **Investment Limit**: Total investment amount cannot exceed the principal amount
3. **Approval Requirements**: Loans cannot transition to approved without approval information
4. **Disbursement Requirements**: Loans cannot transition to disbursed without disbursement information
5. **Unique Investor per Loan**: Each investor can only invest once in the same loan

## Triggers for Business Logic

```sql
-- Trigger to update total_invested_amount when investments are added/updated
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

CREATE TRIGGER trigger_update_loan_total_invested
    AFTER INSERT OR UPDATE ON loan_investments
    FOR EACH ROW
    EXECUTE FUNCTION update_loan_total_invested();

-- Trigger to validate state transitions
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