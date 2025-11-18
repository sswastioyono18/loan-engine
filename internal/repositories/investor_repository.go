package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kitabisa/loan-engine/internal/models"
)

type InvestorRepository interface {
	Create(ctx context.Context, investor *models.Investor) error
	GetByID(ctx context.Context, id int) (*models.Investor, error)
	GetByInvestorID(ctx context.Context, investorID string) (*models.Investor, error)
	GetByEmail(ctx context.Context, email string) (*models.Investor, error)
	Update(ctx context.Context, investor *models.Investor) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, offset, limit int) ([]*models.Investor, error)
}

type investorRepositoryImpl struct {
	base *BaseRepository
}

func NewInvestorRepository(driver Driver) InvestorRepository {
	return &investorRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *investorRepositoryImpl) Create(ctx context.Context, investor *models.Investor) error {
	query := `
		INSERT INTO investors (investor_id, name, email, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.base.GetUtilDB().QueryRowContext(
		ctx, query,
		investor.InvestorID, investor.FullName, investor.Email, investor.Phone,
	).Scan(&investor.ID, &investor.CreatedAt, &investor.UpdatedAt)

	return err
}

func (r *investorRepositoryImpl) GetByID(ctx context.Context, id int) (*models.Investor, error) {
	query := `
		SELECT id, investor_id, name, email, phone, created_at, updated_at
		FROM investors WHERE id = $1
	`

	var investor models.Investor
	err := r.base.GetUtilDB().GetContext(ctx, &investor, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("investor not found")
		}
		return nil, err
	}

	return &investor, nil
}

func (r *investorRepositoryImpl) GetByInvestorID(ctx context.Context, investorID string) (*models.Investor, error) {
	query := `
		SELECT id, investor_id, name, email, phone, created_at, updated_at
		FROM investors WHERE investor_id = $1
	`

	var investor models.Investor
	err := r.base.GetUtilDB().GetContext(ctx, &investor, query, investorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("investor not found")
		}
		return nil, err
	}

	return &investor, nil
}

func (r *investorRepositoryImpl) GetByEmail(ctx context.Context, email string) (*models.Investor, error) {
	query := `
		SELECT id, investor_id, name, email, phone, created_at, updated_at
		FROM investors WHERE email = $1
	`

	var investor models.Investor
	err := r.base.GetUtilDB().GetContext(ctx, &investor, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("investor not found")
		}
		return nil, err
	}

	return &investor, nil
}

func (r *investorRepositoryImpl) Update(ctx context.Context, investor *models.Investor) error {
	query := `
		UPDATE investors SET
			investor_id = $1, name = $2, email = $3,
			phone = $4, updated_at = NOW()
		WHERE id = $5
	`

	result, err := r.base.GetUtilDB().ExecContext(
		ctx, query,
		investor.InvestorID, investor.FullName, investor.Email,
		investor.Phone, investor.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("investor not found")
	}

	return nil
}

func (r *investorRepositoryImpl) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM investors WHERE id = $1"
	result, err := r.base.GetUtilDB().ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("investor not found")
	}

	return nil
}

func (r *investorRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Investor, error) {
	query := `
		SELECT id, investor_id, name, email, phone, created_at, updated_at
		FROM investors
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var investors []*models.Investor
	err := r.base.GetUtilDB().SelectContext(ctx, &investors, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return investors, nil
}