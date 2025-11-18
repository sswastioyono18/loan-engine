package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sswastioyono18/loan-engine/internal/models"
)

type LoanStateHistoryRepository interface {
	Create(ctx context.Context, history *models.LoanStateHistory) error
	GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanStateHistory, error)
	GetLatestByLoanID(ctx context.Context, loanID int) (*models.LoanStateHistory, error)
	List(ctx context.Context, loanID int, offset, limit int) ([]*models.LoanStateHistory, error)
}

type loanStateHistoryRepositoryImpl struct {
	base *BaseRepository
}

func NewLoanStateHistoryRepository(driver Driver) LoanStateHistoryRepository {
	return &loanStateHistoryRepositoryImpl{
		base: NewBaseRepository(driver),
	}
}

func (r *loanStateHistoryRepositoryImpl) Create(ctx context.Context, history *models.LoanStateHistory) error {
	query := `
		INSERT INTO loan_state_history (
			loan_id, old_state, new_state, reason
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.base.GetUtilDB().QueryRowContext(
		ctx, query,
		history.LoanID, history.PreviousState, history.NewState, history.TransitionReason,
	).Scan(&history.ID, &history.CreatedAt)

	return err
}

func (r *loanStateHistoryRepositoryImpl) GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanStateHistory, error) {
	query := `
		SELECT id, loan_id, old_state, new_state, reason, created_at
		FROM loan_state_history WHERE loan_id = $1
		ORDER BY created_at ASC
	`

	var histories []*models.LoanStateHistory
	err := r.base.GetUtilDB().SelectContext(ctx, &histories, query, loanID)
	if err != nil {
		return nil, err
	}

	return histories, nil
}

func (r *loanStateHistoryRepositoryImpl) GetLatestByLoanID(ctx context.Context, loanID int) (*models.LoanStateHistory, error) {
	query := `
		SELECT id, loan_id, old_state, new_state, reason, created_at
		FROM loan_state_history
		WHERE loan_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	var history models.LoanStateHistory
	err := r.base.GetUtilDB().GetContext(ctx, &history, query, loanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no state history found for loan")
		}
		return nil, err
	}

	return &history, nil
}

func (r *loanStateHistoryRepositoryImpl) List(ctx context.Context, loanID int, offset, limit int) ([]*models.LoanStateHistory, error) {
	query := `
		SELECT id, loan_id, old_state, new_state, reason, created_at
		FROM loan_state_history
		WHERE loan_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	var histories []*models.LoanStateHistory
	err := r.base.GetUtilDB().SelectContext(ctx, &histories, query, loanID, limit, offset)
	if err != nil {
		return nil, err
	}

	return histories, nil
}
