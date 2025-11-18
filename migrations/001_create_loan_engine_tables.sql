-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS borrowers (
    id SERIAL PRIMARY KEY,
    id_number VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS investors (
    id SERIAL PRIMARY KEY,
    investor_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    loan_id VARCHAR(50) UNIQUE DEFAULT NULL,
    borrower_id INTEGER NOT NULL,
    principal_amount DECIMAL(15, 2) NOT NULL,
    rate DECIMAL(5, 4) NOT NULL,
    roi DECIMAL(5, 4) NOT NULL,
    agreement_letter_link TEXT,
    current_state VARCHAR(20) DEFAULT 'proposed',
    total_invested_amount DECIMAL(15, 2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (borrower_id) REFERENCES borrowers(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS loan_approvals (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    field_validator_employee_id VARCHAR(50) NOT NULL,
    proof_image_url TEXT,
    approved_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS loan_investments (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    investor_id INTEGER NOT NULL,
    investment_amount DECIMAL(15, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id) ON DELETE CASCADE,
    FOREIGN KEY (investor_id) REFERENCES investors(id) ON DELETE CASCADE,
    UNIQUE(loan_id, investor_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS loan_disbursements (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    field_officer_employee_id VARCHAR(50) NOT NULL,
    agreement_letter_signed_url TEXT,
    disbursed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS loan_state_history (
    id SERIAL PRIMARY KEY,
    loan_id INTEGER NOT NULL,
    old_state VARCHAR(20),
    new_state VARCHAR(20) NOT NULL,
    changed_by VARCHAR(100),
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (loan_id) REFERENCES loans(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_borrowers_id_number ON borrowers(id_number);
CREATE INDEX IF NOT EXISTS idx_investors_email ON investors(email);
CREATE INDEX IF NOT EXISTS idx_loans_borrower_id ON loans(borrower_id);
CREATE INDEX IF NOT EXISTS idx_loans_current_state ON loans(current_state);
CREATE INDEX IF NOT EXISTS idx_loans_loan_id ON loans(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_approvals_loan_id ON loan_approvals(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_investments_loan_id ON loan_investments(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_investments_investor_id ON loan_investments(investor_id);
CREATE INDEX IF NOT EXISTS idx_loan_disbursements_loan_id ON loan_disbursements(loan_id);
CREATE INDEX IF NOT EXISTS idx_loan_state_history_loan_id ON loan_state_history(loan_id);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_borrowers_updated_at ON borrowers;
DROP TRIGGER IF EXISTS update_investors_updated_at ON investors;
DROP TRIGGER IF EXISTS update_loans_updated_at ON loans;
DROP TRIGGER IF EXISTS update_loan_approvals_updated_at ON loan_approvals;
DROP TRIGGER IF EXISTS update_loan_investments_updated_at ON loan_investments;
DROP TRIGGER IF EXISTS update_loan_disbursements_updated_at ON loan_disbursements;
DROP TRIGGER IF EXISTS update_loan_state_history_updated_at ON loan_state_history;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER update_borrowers_updated_at BEFORE UPDATE ON borrowers FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_investors_updated_at BEFORE UPDATE ON investors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_loans_updated_at BEFORE UPDATE ON loans FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_loan_approvals_updated_at BEFORE UPDATE ON loan_approvals FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_loan_investments_updated_at BEFORE UPDATE ON loan_investments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_loan_disbursements_updated_at BEFORE UPDATE ON loan_disbursements FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_loan_state_history_updated_at BEFORE UPDATE ON loan_state_history FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION validate_loan_state_transition()
RETURNS TRIGGER AS $$
BEGIN
    -- Only allow state transitions to move forward
    IF NEW.current_state = 'proposed' THEN
        -- Cannot go back to proposed
        RETURN OLD;
    ELSIF NEW.current_state = 'approved' THEN
        -- Can only transition from proposed
        IF OLD.current_state != 'proposed' THEN
            RAISE EXCEPTION 'Loan can only be approved from proposed state';
        END IF;
    ELSIF NEW.current_state = 'invested' THEN
        -- Can only transition from approved
        IF OLD.current_state != 'approved' THEN
            RAISE EXCEPTION 'Loan can only be invested from approved state';
        END IF;
    ELSIF NEW.current_state = 'disbursed' THEN
        -- Can only transition from invested
        IF OLD.current_state != 'invested' THEN
            RAISE EXCEPTION 'Loan can only be disbursed from invested state';
        END IF;
    END IF;

    -- Insert record into loan_state_history
    INSERT INTO loan_state_history (loan_id, old_state, new_state, changed_by, reason)
    VALUES (NEW.id, OLD.current_state, NEW.current_state, 'system', 'State transition');

    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS validate_loan_state ON loans;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER validate_loan_state BEFORE UPDATE OF current_state ON loans FOR EACH ROW EXECUTE FUNCTION validate_loan_state_transition();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION validate_investment_amount()
RETURNS TRIGGER AS $$
DECLARE
    current_total DECIMAL(15, 2);
    loan_principal DECIMAL(15, 2);
BEGIN
    -- Get the current total invested amount for this loan
    SELECT COALESCE(SUM(investment_amount), 0) INTO current_total
    FROM loan_investments
    WHERE loan_id = NEW.loan_id;

    -- Get the principal amount of the loan
    SELECT principal_amount INTO loan_principal
    FROM loans
    WHERE id = NEW.loan_id;

    -- Check if the new investment would exceed the principal
    IF (current_total + NEW.investment_amount) > loan_principal THEN
        RAISE EXCEPTION 'Investment amount exceeds remaining principal';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS validate_investment ON loan_investments;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER validate_investment BEFORE INSERT OR UPDATE ON loan_investments FOR EACH ROW EXECUTE FUNCTION validate_investment_amount();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_loan_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.loan_id IS NULL THEN
        NEW.loan_id := 'LN' || LPAD(NEW.id::TEXT, 8, '0');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- +goose StatementEnd

-- +goose StatementBegin
DROP TRIGGER IF EXISTS generate_loan_id_trigger ON loans;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER generate_loan_id_trigger BEFORE INSERT ON loans FOR EACH ROW EXECUTE FUNCTION generate_loan_id();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS generate_loan_id_trigger ON loans;
DROP TRIGGER IF EXISTS validate_investment ON loan_investments;
DROP TRIGGER IF EXISTS validate_loan_state ON loans;
DROP TRIGGER IF EXISTS update_borrowers_updated_at ON borrowers;
DROP TRIGGER IF EXISTS update_investors_updated_at ON investors;
DROP TRIGGER IF EXISTS update_loans_updated_at ON loans;
DROP TRIGGER IF EXISTS update_loan_approvals_updated_at ON loan_approvals;
DROP TRIGGER IF EXISTS update_loan_investments_updated_at ON loan_investments;
DROP TRIGGER IF EXISTS update_loan_disbursements_updated_at ON loan_disbursements;
DROP TRIGGER IF EXISTS update_loan_state_history_updated_at ON loan_state_history;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
-- +goose StatementEnd

-- +goose StatementBegin
DROP FUNCTION IF EXISTS generate_loan_id();
DROP FUNCTION IF EXISTS validate_investment_amount();
DROP FUNCTION IF EXISTS validate_loan_state_transition();
DROP FUNCTION IF EXISTS update_updated_at_column();
-- +goose StatementEnd

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_borrowers_id_number;
DROP INDEX IF EXISTS idx_investors_email;
DROP INDEX IF EXISTS idx_loans_borrower_id;
DROP INDEX IF EXISTS idx_loans_current_state;
DROP INDEX IF EXISTS idx_loans_loan_id;
DROP INDEX IF EXISTS idx_loan_approvals_loan_id;
DROP INDEX IF EXISTS idx_loan_investments_loan_id;
DROP INDEX IF EXISTS idx_loan_investments_investor_id;
DROP INDEX IF EXISTS idx_loan_disbursements_loan_id;
DROP INDEX IF EXISTS idx_loan_state_history_loan_id;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS loan_state_history;
DROP TABLE IF EXISTS loan_disbursements;
DROP TABLE IF EXISTS loan_investments;
DROP TABLE IF EXISTS loan_approvals;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS investors;
DROP TABLE IF EXISTS borrowers;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd