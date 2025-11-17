package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kitabisa/loan-engine/internal/models"
)

type LoanApprovalRepository interface {
	Create(ctx context.Context, approval *models.LoanApproval) error
	GetByLoanID(ctx context.Context, loanID int) (*models.LoanApproval, error)
	GetByID(ctx context.Context, id int) (*models.LoanApproval, error)
	Update(ctx context.Context, approval *models.LoanApproval) error
	Delete(ctx context.Context, id int) error
}

type loanApprovalRepositoryImpl struct {
	base *BaseRepository
}

func NewLoanApprovalRepository(driver Driver) LoanApprovalRepository {
	return &loanApprovalRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *loanApprovalRepositoryImpl) Create(ctx context.Context, approval *models.LoanApproval) error {
	query := `
		INSERT INTO loan_approvals (
			loan_id, field_validator_employee_id, approval_date, proof_image_url
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.base.GetUtilDB().QueryRowContext(
		ctx, query,
		approval.LoanID, approval.FieldValidatorEmployeeID,
		approval.ApprovalDate, approval.ProofImageUrl,
	).Scan(&approval.ID, &approval.CreatedAt)

	return err
}

func (r *loanApprovalRepositoryImpl) GetByLoanID(ctx context.Context, loanID int) (*models.LoanApproval, error) {
	query := `
		SELECT id, loan_id, field_validator_employee_id, approval_date,
		       proof_image_url, created_at
		FROM loan_approvals WHERE loan_id = $1
	`

	var approval models.LoanApproval
	err := r.base.GetUtilDB().GetContext(ctx, &approval, query, loanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan approval not found")
		}
		return nil, err
	}

	return &approval, nil
}

func (r *loanApprovalRepositoryImpl) GetByID(ctx context.Context, id int) (*models.LoanApproval, error) {
	query := `
		SELECT id, loan_id, field_validator_employee_id, approval_date,
		       proof_image_url, created_at
		FROM loan_approvals WHERE id = $1
	`

	var approval models.LoanApproval
	err := r.base.GetUtilDB().GetContext(ctx, &approval, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan approval not found")
		}
		return nil, err
	}

	return &approval, nil
}

func (r *loanApprovalRepositoryImpl) Update(ctx context.Context, approval *models.LoanApproval) error {
	query := `
		UPDATE loan_approvals SET
			field_validator_employee_id = $1, approval_date = $2,
			proof_image_url = $3
		WHERE id = $4
	`

	result, err := r.base.GetUtilDB().ExecContext(
		ctx, query,
		approval.FieldValidatorEmployeeID, approval.ApprovalDate,
		approval.ProofImageUrl, approval.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loan approval not found")
	}

	return nil
}

func (r *loanApprovalRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM loan_approvals WHERE id = $1"
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loan approval not found")
	}

	return nil
}