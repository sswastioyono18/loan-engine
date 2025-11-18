package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kitabisa/loan-engine/internal/models"
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
	GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error)
}

type loanRepositoryImpl struct {
	base *BaseRepository
}

func NewLoanRepository(driver Driver) LoanRepository {
	return &loanRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *loanRepositoryImpl) Create(ctx context.Context, loan *models.Loan) error {
	query := `
		INSERT INTO loans (
			borrower_id, principal_amount, rate, roi,
			agreement_letter_link, current_state, total_invested_amount
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.base.GetUtilDB().QueryRowContext(
		ctx, query,
		loan.BorrowerID, loan.PrincipalAmount,
		loan.Rate, loan.ROI, loan.AgreementLetterLink,
		loan.CurrentState, loan.TotalInvestedAmount,
	).Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)

	// After creation, fetch the generated loan_id
	if err == nil {
		fetchQuery := "SELECT loan_id FROM loans WHERE id = $1"
		err = r.base.GetUtilDB().GetContext(ctx, &loan.LoanID, fetchQuery, loan.ID)
	}

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
	err := r.base.GetUtilDB().GetContext(ctx, &loan, query, id)
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
	err := r.base.GetUtilDB().GetContext(ctx, &loan, query, loanID)
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

	result, err := r.base.GetUtilDB().ExecContext(
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
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, id)
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
	paramIndex := 1

	if state != nil {
		query += fmt.Sprintf(" WHERE current_state = $%d", paramIndex)
		args = append(args, *state)
		paramIndex++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	args = append(args, limit, offset)

	var loans []*models.Loan
	err := r.base.GetUtilDB().SelectContext(ctx, &loans, query, args...)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *loanRepositoryImpl) UpdateState(ctx context.Context, id int, newState string) error {
	query := "UPDATE loans SET current_state = $1, updated_at = NOW() WHERE id = $2"
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, newState, id)
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
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, amount, loanID)
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
	err := r.base.GetUtilDB().SelectContext(ctx, &loans, query, state)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *loanRepositoryImpl) GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error) {
	query := "SELECT total_invested_amount FROM loans WHERE id = $1"

	var amount float64
	err := r.base.GetUtilDB().GetContext(ctx, &amount, query, loanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("loan not found")
		}
		return 0, err
	}

	return amount, nil
}