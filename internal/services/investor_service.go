package services

import (
	"context"

	"github.com/sswastioyono18/loan-engine/internal/models"
)

type InvestorService interface {
	CreateInvestor(ctx context.Context, investor *models.Investor) error
	GetInvestorByID(ctx context.Context, id int) (*models.Investor, error)
	GetInvestorByInvestorID(ctx context.Context, investorID string) (*models.Investor, error)
	GetInvestorByEmail(ctx context.Context, email string) (*models.Investor, error)
	UpdateInvestor(ctx context.Context, id int, investor *models.Investor) error
	DeleteInvestor(ctx context.Context, id int) error
	ListInvestors(ctx context.Context, offset, limit int) ([]*models.Investor, error)
}

type investorServiceImpl struct {
	repo InvestorRepository
}

func NewInvestorService(repo InvestorRepository) InvestorService {
	return &investorServiceImpl{
		repo: repo,
	}
}

func (s *investorServiceImpl) CreateInvestor(ctx context.Context, investor *models.Investor) error {
	return s.repo.Create(ctx, investor)
}

func (s *investorServiceImpl) GetInvestorByID(ctx context.Context, id int) (*models.Investor, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *investorServiceImpl) GetInvestorByInvestorID(ctx context.Context, investorID string) (*models.Investor, error) {
	return s.repo.GetByInvestorID(ctx, investorID)
}

func (s *investorServiceImpl) GetInvestorByEmail(ctx context.Context, email string) (*models.Investor, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *investorServiceImpl) UpdateInvestor(ctx context.Context, id int, investor *models.Investor) error {
	// Get existing investor to check if it exists
	existingInvestor, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	investor.ID = id
	investor.CreatedAt = existingInvestor.CreatedAt

	return s.repo.Update(ctx, investor)
}

func (s *investorServiceImpl) DeleteInvestor(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *investorServiceImpl) ListInvestors(ctx context.Context, offset, limit int) ([]*models.Investor, error) {
	return s.repo.List(ctx, offset, limit)
}
