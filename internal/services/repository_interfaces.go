package services

import (
	"context"

	"github.com/sswastioyono18/loan-engine/internal/models"
)

// BorrowerRepository defines the specific methods that BorrowerService needs from the repository
type BorrowerRepository interface {
	Create(ctx context.Context, borrower *models.Borrower) error
	GetByID(ctx context.Context, id int) (*models.Borrower, error)
	GetByBorrowerIDNumber(ctx context.Context, borrowerIDNumber string) (*models.Borrower, error)
	Update(ctx context.Context, borrower *models.Borrower) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*models.Borrower, error)
}

// InvestorRepository defines the specific methods that InvestorService and other services need from the investor repository
type InvestorRepository interface {
	GetByID(ctx context.Context, id int) (*models.Investor, error)
	GetByInvestorID(ctx context.Context, investorID string) (*models.Investor, error)
	GetByEmail(ctx context.Context, email string) (*models.Investor, error)
	Create(ctx context.Context, investor *models.Investor) error
	Update(ctx context.Context, investor *models.Investor) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*models.Investor, error)
}

// UserRepository defines the specific methods that AuthService needs from the repository
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUserID(ctx context.Context, userID string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	UpdatePassword(ctx context.Context, id int, hashedPassword string) error
}

// LoanRepository defines the specific methods that LoanService needs from the loan repository
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
	GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error)
}

// LoanApprovalRepository defines the specific methods that LoanService needs from the loan approval repository
type LoanApprovalRepository interface {
	Create(ctx context.Context, approval *models.LoanApproval) error
	GetByLoanID(ctx context.Context, loanID int) (*models.LoanApproval, error)
}

// LoanDisbursementRepository defines the specific methods that LoanService needs from the loan disbursement repository
type LoanDisbursementRepository interface {
	Create(ctx context.Context, disbursement *models.LoanDisbursement) error
	GetByLoanID(ctx context.Context, loanID int) (*models.LoanDisbursement, error)
}

// LoanInvestmentRepository defines the specific methods that LoanService needs from the loan investment repository
type LoanInvestmentRepository interface {
	Create(ctx context.Context, investment *models.LoanInvestment) error
	GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanInvestment, error)
	GetByLoanAndInvestor(ctx context.Context, loanID int, investorID int) (*models.LoanInvestment, error)
}

// LoanStateHistoryRepository defines the specific methods that LoanService needs from the loan state history repository
type LoanStateHistoryRepository interface {
	Create(ctx context.Context, history *models.LoanStateHistory) error
	GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanStateHistory, error)
}
