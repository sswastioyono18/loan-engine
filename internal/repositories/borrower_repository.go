package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kitabisa/loan-engine/internal/models"
)

type BorrowerRepository interface {
	Create(ctx context.Context, borrower *models.Borrower) error
	GetByID(ctx context.Context, id int) (*models.Borrower, error)
	GetByBorrowerIDNumber(ctx context.Context, borrowerIDNumber string) (*models.Borrower, error)
	Update(ctx context.Context, borrower *models.Borrower) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*models.Borrower, error)
}

type borrowerRepositoryImpl struct {
	base *BaseRepository
}

func NewBorrowerRepository(driver Driver) BorrowerRepository {
	return &borrowerRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *borrowerRepositoryImpl) Create(ctx context.Context, borrower *models.Borrower) error {
	query := `
		INSERT INTO borrowers (id_number, name, email, phone, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.base.GetUtilDB().QueryRowContext(
		ctx, query,
		borrower.BorrowerIDNumber, borrower.FullName, borrower.Email,
		borrower.Phone, borrower.Address,
	).Scan(&borrower.ID, &borrower.CreatedAt, &borrower.UpdatedAt)

	return err
}

func (r *borrowerRepositoryImpl) GetByID(ctx context.Context, id int) (*models.Borrower, error) {
	query := `
		SELECT id, id_number, name, email, phone, address, created_at, updated_at
		FROM borrowers WHERE id = $1
	`

	var borrower models.Borrower
	err := r.base.GetUtilDB().GetContext(ctx, &borrower, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("borrower not found")
		}
		return nil, err
	}

	return &borrower, nil
}

func (r *borrowerRepositoryImpl) GetByBorrowerIDNumber(ctx context.Context, borrowerIDNumber string) (*models.Borrower, error) {
	query := `
		SELECT id, id_number, name, email, phone, address, created_at, updated_at
		FROM borrowers WHERE id_number = $1
	`

	var borrower models.Borrower
	err := r.base.GetUtilDB().GetContext(ctx, &borrower, query, borrowerIDNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("borrower not found")
		}
		return nil, err
	}

	return &borrower, nil
}

func (r *borrowerRepositoryImpl) Update(ctx context.Context, borrower *models.Borrower) error {
	query := `
		UPDATE borrowers SET
			id_number = $1, name = $2, email = $3,
			phone = $4, address = $5, updated_at = NOW()
		WHERE id = $6
	`

	result, err := r.base.GetUtilDB().ExecContext(
		ctx, query,
		borrower.BorrowerIDNumber, borrower.FullName, borrower.Email,
		borrower.Phone, borrower.Address, borrower.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("borrower not found")
	}

	return nil
}

func (r *borrowerRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM borrowers WHERE id = $1"
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("borrower not found")
	}

	return nil
}

func (r *borrowerRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Borrower, error) {
	query := `
		SELECT id, id_number, name, email, phone, address, created_at, updated_at
		FROM borrowers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var borrowers []*models.Borrower
	err := r.base.GetUtilDB().SelectContext(ctx, &borrowers, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return borrowers, nil
}
