package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kitabisa/loan-engine/internal/models"
)

type LoanDisbursementRepository interface {
	Create(ctx context.Context, disbursement *models.LoanDisbursement) error
	GetByLoanID(ctx context.Context, loanID int) (*models.LoanDisbursement, error)
	GetByID(ctx context.Context, id int) (*models.LoanDisbursement, error)
	Update(ctx context.Context, disbursement *models.LoanDisbursement) error
	Delete(ctx context.Context, id int) error
}

type loanDisbursementRepositoryImpl struct {
	base *BaseRepository
}

func NewLoanDisbursementRepository(driver Driver) LoanDisbursementRepository {
	return &loanDisbursementRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *loanDisbursementRepositoryImpl) Create(ctx context.Context, disbursement *models.LoanDisbursement) error {
	query := `
		INSERT INTO loan_disbursements (
			loan_id, field_officer_employee_id, disbursement_date,
			agreement_letter_signed_url
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.base.GetUtilDB().QueryRowContext(
		ctx, query,
		disbursement.LoanID, disbursement.FieldOfficerEmployeeID,
		disbursement.DisbursementDate, disbursement.AgreementLetterSignedUrl,
	).Scan(&disbursement.ID, &disbursement.CreatedAt)

	return err
}

func (r *loanDisbursementRepositoryImpl) GetByLoanID(ctx context.Context, loanID int) (*models.LoanDisbursement, error) {
	query := `
		SELECT id, loan_id, field_officer_employee_id, disbursement_date,
		       agreement_letter_signed_url, created_at
		FROM loan_disbursements WHERE loan_id = $1
	`

	var disbursement models.LoanDisbursement
	err := r.base.GetUtilDB().GetContext(ctx, &disbursement, query, loanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan disbursement not found")
		}
		return nil, err
	}

	return &disbursement, nil
}

func (r *loanDisbursementRepositoryImpl) GetByID(ctx context.Context, id int) (*models.LoanDisbursement, error) {
	query := `
		SELECT id, loan_id, field_officer_employee_id, disbursement_date,
		       agreement_letter_signed_url, created_at
		FROM loan_disbursements WHERE id = $1
	`

	var disbursement models.LoanDisbursement
	err := r.base.GetUtilDB().GetContext(ctx, &disbursement, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("loan disbursement not found")
		}
		return nil, err
	}

	return &disbursement, nil
}

func (r *loanDisbursementRepositoryImpl) Update(ctx context.Context, disbursement *models.LoanDisbursement) error {
	query := `
		UPDATE loan_disbursements SET
			field_officer_employee_id = $1, disbursement_date = $2,
			agreement_letter_signed_url = $3
		WHERE id = $4
	`

	result, err := r.base.GetUtilDB().ExecContext(
		ctx, query,
		disbursement.FieldOfficerEmployeeID, disbursement.DisbursementDate,
		disbursement.AgreementLetterSignedUrl, disbursement.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loan disbursement not found")
	}

	return nil
}

func (r *loanDisbursementRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM loan_disbursements WHERE id = $1"
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("loan disbursement not found")
	}

	return nil
}