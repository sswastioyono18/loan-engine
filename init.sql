-- init.sql - Database initialization script for loan engine

-- Create the database schema
CREATE SCHEMA IF NOT EXISTS loan_engine;

-- Set the search path to use the loan_engine schema
SET search_path TO loan_engine, public;

-- Create borrowers table
CREATE TABLE IF NOT EXISTS borrowers (
    id SERIAL PRIMARY KEY,
    borrower_id_number VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create users table (for employees/validators)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    employee_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL, -- 'admin', 'validator', 'field_officer', 'investor'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loans table
CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    loan_id VARCHAR(50) UNIQUE NOT NULL,
    borrower_id INTEGER NOT NULL REFERENCES borrowers(id) ON DELETE CASCADE,
    principal_amount DECIMAL(15, 2) NOT NULL,
    rate DECIMAL(5, 2) NOT NULL, -- interest rate percentage
    roi DECIMAL(5, 2) NOT NULL, -- return of investment percentage
    agreement_letter_url TEXT,
    state VARCHAR(20) NOT NULL DEFAULT 'proposed', -- 'proposed', 'approved', 'invested', 'disbursed'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (state IN ('proposed', 'approved', 'invested', 'disbursed'))
);

-- Create loan_state_history table to track state changes
CREATE TABLE IF NOT EXISTS loan_state_history (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    old_state VARCHAR(20),
    new_state VARCHAR(20) NOT NULL,
    changed_by INTEGER REFERENCES users(id), -- user who changed the state
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);

-- Create loan_approvals table
CREATE TABLE IF NOT EXISTS loan_approvals (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    field_validator_employee_id VARCHAR(50) NOT NULL,
    proof_image_url TEXT NOT NULL,
    approved_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create investors table
CREATE TABLE IF NOT EXISTS investors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create loan_investments table
CREATE TABLE IF NOT EXISTS loan_investments (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    investor_id INTEGER NOT NULL REFERENCES investors(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(loan_id, investor_id) -- Each investor can invest in a loan only once
);

-- Create loan_disbursements table
CREATE TABLE IF NOT EXISTS loan_disbursements (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    agreement_letter_url TEXT NOT NULL,
    field_officer_employee_id VARCHAR(50) NOT NULL,
    disbursement_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger function to enforce state transition rules
CREATE OR REPLACE FUNCTION enforce_state_transition()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the state is changing
    IF NEW.state <> OLD.state THEN
        -- Ensure state only moves forward
        CASE
            WHEN OLD.state = 'proposed' THEN
                -- From proposed, can only go to approved
                IF NEW.state NOT IN ('approved') THEN
                    RAISE EXCEPTION 'Invalid state transition: % to %', OLD.state, NEW.state;
                END IF;
            WHEN OLD.state = 'approved' THEN
                -- From approved, can go to invested or stay approved
                IF NEW.state NOT IN ('invested') THEN
                    RAISE EXCEPTION 'Invalid state transition: % to %', OLD.state, NEW.state;
                END IF;
            WHEN OLD.state = 'invested' THEN
                -- From invested, can go to disbursed or stay invested
                IF NEW.state NOT IN ('disbursed') THEN
                    RAISE EXCEPTION 'Invalid state transition: % to %', OLD.state, NEW.state;
                END IF;
            WHEN OLD.state = 'disbursed' THEN
                -- From disbursed, no further transitions allowed
                RAISE EXCEPTION 'Invalid state transition: % to %', OLD.state, NEW.state;
        END CASE;

        -- Insert record into loan_state_history
        INSERT INTO loan_state_history (loan_id, old_state, new_state, changed_by, notes)
        VALUES (NEW.id, OLD.state, NEW.state, NEW.updated_by, 'State transition');
    END IF;

    -- Update the updated_at timestamp
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to enforce state transition rules
CREATE TRIGGER loan_state_transition_trigger
    BEFORE UPDATE ON loans
    FOR EACH ROW
    EXECUTE FUNCTION enforce_state_transition();

-- Create trigger function to validate investment amount
CREATE OR REPLACE FUNCTION validate_investment_amount()
RETURNS TRIGGER AS $$
DECLARE
    total_invested DECIMAL(15, 2);
    loan_principal DECIMAL(15, 2);
BEGIN
    -- Get the principal amount of the loan
    SELECT principal_amount INTO loan_principal
    FROM loans
    WHERE id = NEW.loan_id;

    -- Calculate total invested amount for this loan (including the new investment)
    SELECT COALESCE(SUM(amount), 0) INTO total_invested
    FROM loan_investments
    WHERE loan_id = NEW.loan_id;

    -- Add the new investment amount
    total_invested := total_invested + NEW.amount;

    -- Check if total invested exceeds the principal amount
    IF total_invested > loan_principal THEN
        RAISE EXCEPTION 'Total investment amount exceeds loan principal amount';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to validate investment amount
CREATE TRIGGER investment_amount_validation_trigger
    BEFORE INSERT OR UPDATE ON loan_investments
    FOR EACH ROW
    EXECUTE FUNCTION validate_investment_amount();

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_borrowers_id_number ON borrowers(borrower_id_number);
CREATE INDEX IF NOT EXISTS idx_loans_state ON loans(state);
CREATE INDEX IF NOT EXISTS idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX IF NOT EXISTS idx_loan_approvals_loan_id ON loan_approvals(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_investments_loan_id ON loan_investments(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_investments_investor_id ON loan_investments(investor_id);
CREATE INDEX IF NOT EXISTS idx_loan_disbursements_loan_id ON loan_disbursements(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_state_history_loan_id ON loan_state_history(loan_id);

-- Insert sample data for testing
INSERT INTO borrowers (borrower_id_number, name, email, phone, address) VALUES
('B001', 'John Doe', 'john.doe@example.com', '1234567890', '123 Main St, City, Country'),
('B002', 'Jane Smith', 'jane.smith@example.com', '0987654321', '456 Oak Ave, City, Country');

INSERT INTO users (employee_id, name, email, role) VALUES
('E001', 'Alice Johnson', 'alice.johnson@example.com', 'validator'),
('E002', 'Bob Williams', 'bob.williams@example.com', 'field_officer'),
('E003', 'Carol Brown', 'carol.brown@example.com', 'admin');

INSERT INTO investors (name, email, phone) VALUES
('Investor 1', 'investor1@example.com', '1111111111'),
('Investor 2', 'investor2@example.com', '2222222222'),
('Investor 3', 'investor3@example.com', '3333333333');