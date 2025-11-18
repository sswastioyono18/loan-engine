package services

import (
	"context"
	"github.com/sswastioyono18/loan-engine/internal/models"
	"github.com/sswastioyono18/loan-engine/internal/repositories"
)

type BorrowerService interface {
	CreateBorrower(ctx context.Context, borrower *models.Borrower) error
	GetBorrowerByID(ctx context.Context, id int) (*models.Borrower, error)
	GetBorrowerByBorrowerIDNumber(ctx context.Context, borrowerIDNumber string) (*models.Borrower, error)
	UpdateBorrower(ctx context.Context, id int, borrower *models.Borrower) error
	DeleteBorrower(ctx context.Context, id int) error
	ListBorrowers(ctx context.Context, offset, limit int) ([]*models.Borrower, error)
}

type borrowerServiceImpl struct {
	repo repositories.BorrowerRepository
}

func NewBorrowerService(repo repositories.BorrowerRepository) BorrowerService {
	return &borrowerServiceImpl{
		repo: repo,
	}
}

func (s *borrowerServiceImpl) CreateBorrower(ctx context.Context, borrower *models.Borrower) error {
	return s.repo.Create(ctx, borrower)
}

func (s *borrowerServiceImpl) GetBorrowerByID(ctx context.Context, id int) (*models.Borrower, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *borrowerServiceImpl) GetBorrowerByBorrowerIDNumber(ctx context.Context, borrowerIDNumber string) (*models.Borrower, error) {
	return s.repo.GetByBorrowerIDNumber(ctx, borrowerIDNumber)
}

func (s *borrowerServiceImpl) UpdateBorrower(ctx context.Context, id int, borrower *models.Borrower) error {
	// Get existing borrower to check if it exists
	existingBorrower, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	borrower.ID = id
	borrower.CreatedAt = existingBorrower.CreatedAt

	return s.repo.Update(ctx, borrower)
}

func (s *borrowerServiceImpl) DeleteBorrower(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *borrowerServiceImpl) ListBorrowers(ctx context.Context, offset, limit int) ([]*models.Borrower, error) {
	return s.repo.List(ctx, offset, limit)
}
