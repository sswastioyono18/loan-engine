# Loan Engine Requirements Analysis

## Overview
This document provides a comprehensive analysis of how the loan engine system satisfies all the original requirements specified for the loan lifecycle management system. The system supports loans through four states: proposed, approved, invested, and disbursed.

## Original Requirements Summary
The loan engine was designed to support a loan lifecycle with the following states and rules:
1. **Proposed** - Initial state when loan is created
2. **Approved** - After staff approval with required documentation
3. **Invested** - After sufficient investment from investors
4. **Disbursed** - After loan is given to borrower

State transitions can only move forward, never backward.

## Requirements Implementation Analysis

### 1. Proposed State (Initial State)
**Requirement**: Proposed is the initial state when loan is created

**Implementation**:
- In [`internal/services/loan_service.go:80`](internal/services/loan_service.go:80), loans are created with `loan.CurrentState = "proposed"` as the default state
- Code: `loan.CurrentState = "proposed"` in the `CreateLoan` method
- Database schema in [`migrations/001_create_loan_engine_tables.sql:36`](migrations/001_create_loan_engine_tables.sql:36) sets default state to 'proposed'

### 2. Approved State Requirements
**Requirement**: Approval must contain picture proof of field validator visit, employee ID, and date of approval

**Implementation**:
- The [`internal/models/loan_approval.go`](internal/models/loan_approval.go) model contains:
  - `ProofImageUrl` for picture proof of field validator visit
  - `FieldValidatorEmployeeID` for employee ID of field validator
  - `ApprovalDate` for date of approval
- API endpoint: `POST /api/v1/loans/{id}/approve` in [`internal/handlers/loan_handler.go:157-184`](internal/handlers/loan_handler.go:157-184)
- Validation ensures all required approval information is provided

### 3. State Transition Rules
**Requirement**: Once approved, loan cannot go back to proposed state

**Implementation**: Multiple layers of validation:
- Database trigger in [`migrations/001_create_loan_engine_tables.sql:156-194`](migrations/001_create_loan_engine_tables.sql:156-194) prevents backward transitions
- Service layer validation in [`internal/services/loan_service.go:148-150`](internal/services/loan_service.go:148-150) checks `if loan.CurrentState != "proposed"`
- Documentation in [`docs/state_transition_validation.md:160-179`](docs/state_transition_validation.md:160-179)
- Once approved, loan is ready to be offered to investors/lenders

### 4. Investment Requirements
**Requirement**: Invested state when total investment equals principal amount

**Implementation**:
- In [`internal/services/loan_service.go:234-239`](internal/services/loan_service.go:234-239), when `newTotal >= loan.PrincipalAmount`, the loan state is updated to "invested"
- **Multiple Investors**: The [`internal/models/loan_investment.go`](internal/models/loan_investment.go) model supports multiple investors with different amounts via the `InvestorID` and `InvestmentAmount` fields
- **Principal Limit**: In [`internal/services/loan_service.go:207-211`](internal/services/loan_service.go:207-211), there's validation to ensure `investment.InvestmentAmount > remainingPrincipal` doesn't exceed the remaining principal
- Database constraint in [`migrations/001_create_loan_engine_tables.sql:229`](migrations/001_create_loan_engine_tables.sql:229) prevents investment amount from exceeding remaining principal

### 5. Email Notifications
**Requirement**: Investors receive email with agreement letter link when loan is fully invested

**Implementation**:
- In [`internal/services/loan_service.go:254-275`](internal/services/loan_service.go:254-275), when a loan becomes fully invested, the system sends investment confirmation emails to all investors with the agreement letter link
- Email service interface in [`pkg/external/email_service.go`](pkg/external/email_service.go) handles sending notifications
- Agreement letter link is included in the email to investors

### 6. Disbursement Requirements
**Requirement**: Disbursement must contain signed agreement letter, employee ID, and date

**Implementation**:
- The [`internal/models/loan_disbursement.go`](internal/models/loan_disbursement.go) model contains:
  - `AgreementLetterSignedUrl` for signed agreement letter (PDF/JPEG)
  - `FieldOfficerEmployeeID` for employee ID of field officer
  - `DisbursementDate` for date of disbursement
- API endpoint: `POST /api/v1/loans/{id}/disburse` in [`internal/handlers/loan_handler.go:217-244`](internal/handlers/loan_handler.go:217-244)

### 7. Forward-Only State Transitions
**Requirement**: Movement between states can only move forward

**Implementation**: Multiple validation layers:
- Database triggers in [`migrations/001_create_loan_engine_tables.sql:156-194`](migrations/001_create_loan_engine_tables.sql:156-194)
- Service layer validation in [`internal/services/loan_service.go:340-368`](internal/services/loan_service.go:340-368) with `CanTransitionToState` method
- Documentation in [`docs/state_transition_validation.md`](docs/state_transition_validation.md)

### 8. Required Loan Information
**Requirement**: Loan must contain borrower ID number, principal amount, rate, ROI, and agreement letter link

**Implementation**: The [`internal/models/loan.go`](internal/models/loan.go) model contains:
- `BorrowerID` for borrower identification
- `PrincipalAmount` for loan principal
- `Rate` for interest rate (will define total interest that borrower will pay)
- `ROI` for return of investment (will define total profit received by investors)
- `AgreementLetterLink` for link to the generated agreement letter

### 9. RESTful API Implementation
**Requirement**: Design a RESTful API that satisfies the above requirements

**Implementation**: Complete API implementation with:
- `POST /api/v1/loans` - Create loan (proposed state)
- `POST /api/v1/loans/{id}/approve` - Approve loan
- `POST /api/v1/loans/{id}/invest` - Invest in loan
- `POST /api/v1/loans/{id}/disburse` - Disburse loan
- Complete CRUD operations for loans, borrowers, and investors
- All endpoints documented in [`docs/api_endpoints.md`](docs/api_endpoints.md)

### 10. Additional Features Implemented
- **State History Tracking**: Complete audit trail with [`internal/models/loan_state_history.go`](internal/models/loan_state_history.go)
- **Investor Management**: Full investor CRUD operations
- **Email Service**: Mock email service for notifications
- **File Storage**: Mock storage service for document handling
- **Comprehensive Testing**: Unit tests for all major functionality

## System Architecture

### Technology Stack
- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Router**: Chi router
- **Database Operations**: SQLX
- **Pattern**: Repository and Service Layer patterns

### Architecture Components
1. **Models**: Data structures in `internal/models/`
2. **Repositories**: Data access layer in `internal/repositories/`
3. **Services**: Business logic in `internal/services/`
4. **Handlers**: HTTP endpoints in `internal/handlers/`
5. **External Services**: Email and storage services in `pkg/external/`

## Limitation & assumption
- Assumed that same borrower can submit multiple loan request. I did not validate any of this
- System did not validate whether inside the picture proof validator visited the borrower is real. It just checks whether the payload is empty or not
- System did not validate the content signed loan agreement process. I assume this is handled outside of the system. The system just knew the url is there. There is no logic to check is it real signed or not.
- Sending email is mocked
- No auth, RBAC or token implemented for now to request, approve and disbursed
- No caching implemented yet (redis is there though)
- No observability
- Technical implementation of json and db shared a struct. I assume this is a PoC, i dont really care about splitting entity/domain model
