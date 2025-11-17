package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kitabisa/loan-engine/internal/models"
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

type loanInvestmentRepositoryImpl struct {
	base *BaseRepository
}

func NewLoanInvestmentRepository(driver Driver) LoanInvestmentRepository {
	return &loanInvestmentRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *loanInvestmentRepositoryImpl) Create(ctx context.Context, investment *models.LoanInvestment) error {
	query := `
		INSERT INTO loan_investments (loan_id, investor_id, investment_amount)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.base.GetUtilDB().QueryRowContext(
		ctx, query,
		investment.LoanID, investment.InvestorID, investment.InvestmentAmount,
	).Scan(&investment.ID, &investment.CreatedAt)

	return err
}

func (r *loanInvestmentRepositoryImpl) GetByID(ctx context.Context, id int) (*models.LoanInvestment, error) {
	query := `
		SELECT id, loan_id, investor_id, investment_amount, created_at
		FROM loan_investments WHERE id = $1
	`

	var investment models.LoanInvestment
	err := r.base.GetUtilDB().GetContext(ctx, &investment, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan investment not found")
		}
		return nil, err
	}

	return &investment, nil
}

func (r *loanInvestmentRepositoryImpl) GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanInvestment, error) {
	query := `
		SELECT id, loan_id, investor_id, investment_amount, created_at
		FROM loan_investments WHERE loan_id = $1
		ORDER BY created_at DESC
	`

	var investments []*models.LoanInvestment
	err := r.base.GetUtilDB().SelectContext(ctx, &investments, query, loanID)
	if err != nil {
		return nil, err
	}

	return investments, nil
}

func (r *loanInvestmentRepositoryImpl) GetByInvestorID(ctx context.Context, investorID int) ([]*models.LoanInvestment, error) {
	query := `
		SELECT id, loan_id, investor_id, investment_amount, created_at
		FROM loan_investments WHERE investor_id = $1
		ORDER BY created_at DESC
	`

	var investments []*models.LoanInvestment
	err := r.base.GetUtilDB().SelectContext(ctx, &investments, query, investorID)
	if err != nil {
		return nil, err
	}

	return investments, nil
}

func (r *loanInvestmentRepositoryImpl) GetByLoanAndInvestor(ctx context.Context, loanID, investorID int) (*models.LoanInvestment, error) {
	query := `
		SELECT id, loan_id, investor_id, investment_amount, created_at
		FROM loan_investments WHERE loan_id = $1 AND investor_id = $2
	`

	var investment models.LoanInvestment
	err := r.base.GetUtilDB().GetContext(ctx, &investment, query, loanID, investorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan investment not found")
		}
		return nil, err
	}

	return &investment, nil
}

func (r *loanInvestmentRepositoryImpl) Update(ctx context.Context, investment *models.LoanInvestment) error {
	query := `
		UPDATE loan_investments SET
			investment_amount = $1
		WHERE id = $2
	`

	result, err := r.base.GetUtilDB().ExecContext(
		ctx, query,
		investment.InvestmentAmount, investment.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loan investment not found")
	}

	return nil
}

func (r *loanInvestmentRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM loan_investments WHERE id = $1"
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loan investment not found")
	}

	return nil
}

func (r *loanInvestmentRepositoryImpl) GetTotalInvestedAmountByLoan(ctx context.Context, loanID int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(investment_amount), 0)
		FROM loan_investments
		WHERE loan_id = $1
	`

	var total float64
	err := r.base.GetUtilDB().GetContext(ctx, &total, query, loanID)
	if err != nil {
		return 0, err
	}

	return total, nil
}