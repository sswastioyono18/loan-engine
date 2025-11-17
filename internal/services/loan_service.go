package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/kitabisa/loan-engine/internal/models"
	"github.com/kitabisa/loan-engine/internal/repositories"
	"github.com/kitabisa/loan-engine/pkg/external"
)

type LoanService interface {
	CreateLoan(ctx context.Context, loan *models.Loan) error
	GetLoanByID(ctx context.Context, id int) (*models.Loan, error)
	GetLoanByLoanID(ctx context.Context, loanID string) (*models.Loan, error)
	UpdateLoan(ctx context.Context, id int, loan *models.Loan) error
	DeleteLoan(ctx context.Context, id int) error
	ListLoans(ctx context.Context, state *string, offset, limit int) ([]*models.Loan, error)
	GetLoansByState(ctx context.Context, state string) ([]*models.Loan, error)
	
	// State transition methods
	ApproveLoan(ctx context.Context, loanID int, approvalData *models.LoanApproval) error
	InvestInLoan(ctx context.Context, loanID int, investment *models.LoanInvestment) error
	DisburseLoan(ctx context.Context, loanID int, disbursementData *models.LoanDisbursement) error
	
	// Helper methods
	GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error)
	CanTransitionToState(ctx context.Context, loanID int, newState string) (bool, error)
}

type loanServiceImpl struct {
	loanRepo             repositories.LoanRepository
	loanApprovalRepo     repositories.LoanApprovalRepository
	loanDisbursementRepo repositories.LoanDisbursementRepository
	loanInvestmentRepo   repositories.LoanInvestmentRepository
	loanStateHistoryRepo repositories.LoanStateHistoryRepository
	investorRepo         repositories.InvestorRepository
	emailService         external.EmailService
	storageService       external.StorageService
}

func NewLoanService(
	loanRepo repositories.LoanRepository,
	loanApprovalRepo repositories.LoanApprovalRepository,
	loanDisbursementRepo repositories.LoanDisbursementRepository,
	loanInvestmentRepo repositories.LoanInvestmentRepository,
	loanStateHistoryRepo repositories.LoanStateHistoryRepository,
	investorRepo repositories.InvestorRepository,
	emailService external.EmailService,
	storageService external.StorageService,
) LoanService {
	return &loanServiceImpl{
		loanRepo:             loanRepo,
		loanApprovalRepo:     loanApprovalRepo,
		loanDisbursementRepo: loanDisbursementRepo,
		loanInvestmentRepo:   loanInvestmentRepo,
		loanStateHistoryRepo: loanStateHistoryRepo,
		investorRepo:         investorRepo,
		emailService:         emailService,
		storageService:       storageService,
	}
}

func (s *loanServiceImpl) CreateLoan(ctx context.Context, loan *models.Loan) error {
	// Validate required fields
	if loan.PrincipalAmount <= 0 {
		return errors.New("principal amount must be greater than 0")
	}
	
	if loan.Rate < 0 || loan.Rate > 100 {
		return errors.New("rate must be between 0 and 100")
	}
	
	if loan.ROI < 0 || loan.ROI > 100 {
		return errors.New("ROI must be between 0 and 100")
	}
	
	// Set initial state to proposed
	loan.CurrentState = "proposed"
	loan.TotalInvestedAmount = 0.0
	
	return s.loanRepo.Create(ctx, loan)
}

func (s *loanServiceImpl) GetLoanByID(ctx context.Context, id int) (*models.Loan, error) {
	return s.loanRepo.GetByID(ctx, id)
}

func (s *loanServiceImpl) GetLoanByLoanID(ctx context.Context, loanID string) (*models.Loan, error) {
	return s.loanRepo.GetByLoanID(ctx, loanID)
}

func (s *loanServiceImpl) UpdateLoan(ctx context.Context, id int, loan *models.Loan) error {
	// Get existing loan to check state
	existingLoan, err := s.loanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Prevent modification of certain fields based on state
	if existingLoan.CurrentState != "proposed" {
		// Only allow updating specific fields after loan is approved
		loan.BorrowerID = existingLoan.BorrowerID
		loan.PrincipalAmount = existingLoan.PrincipalAmount
		loan.Rate = existingLoan.Rate
		loan.ROI = existingLoan.ROI
		loan.AgreementLetterLink = existingLoan.AgreementLetterLink
	}
	
	// Update fields
	loan.ID = id
	loan.CreatedAt = existingLoan.CreatedAt
	
	return s.loanRepo.Update(ctx, loan)
}

func (s *loanServiceImpl) DeleteLoan(ctx context.Context, id int) error {
	// Check if loan can be deleted (must be in proposed state)
	loan, err := s.loanRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if loan.CurrentState != "proposed" {
		return errors.New("loan can only be deleted in proposed state")
	}
	
	return s.loanRepo.Delete(ctx, id)
}

func (s *loanServiceImpl) ListLoans(ctx context.Context, state *string, offset, limit int) ([]*models.Loan, error) {
	return s.loanRepo.List(ctx, state, offset, limit)
}

func (s *loanServiceImpl) GetLoansByState(ctx context.Context, state string) ([]*models.Loan, error) {
	return s.loanRepo.GetByState(ctx, state)
}

func (s *loanServiceImpl) ApproveLoan(ctx context.Context, loanID int, approvalData *models.LoanApproval) error {
	// Get the loan
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}
	
	// Check if loan is in proposed state
	if loan.CurrentState != "proposed" {
		return errors.New("loan must be in proposed state to be approved")
	}
	
	// Validate approval data
	if approvalData.FieldValidatorEmployeeID == "" {
		return errors.New("field validator employee ID is required")
	}
	
	if approvalData.ProofImageUrl == "" {
		return errors.New("proof image URL is required")
	}
	
	// Create loan approval record
	approvalData.LoanID = loanID
	err = s.loanApprovalRepo.Create(ctx, approvalData)
	if err != nil {
		return fmt.Errorf("failed to create loan approval: %w", err)
	}
	
	// Update loan state to approved
	err = s.loanRepo.UpdateState(ctx, loanID, "approved")
	if err != nil {
		return fmt.Errorf("failed to update loan state: %w", err)
	}
	
	// Add state transition to history
	stateHistory := &models.LoanStateHistory{
		LoanID:           loanID,
		PreviousState:    loan.CurrentState,
		NewState:         "approved",
		TransitionReason: "Loan approved by staff",
	}
	
	err = s.loanStateHistoryRepo.Create(ctx, stateHistory)
	if err != nil {
		return fmt.Errorf("failed to create state history: %w", err)
	}
	
	return nil
}

func (s *loanServiceImpl) InvestInLoan(ctx context.Context, loanID int, investment *models.LoanInvestment) error {
	// Get the loan
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}
	
	// Check if loan is in approved state
	if loan.CurrentState != "approved" {
		return errors.New("loan must be in approved state to receive investments")
	}
	
	// Validate investment amount
	if investment.InvestmentAmount <= 0 {
		return errors.New("investment amount must be greater than 0")
	}
	
	// Check if investment amount exceeds remaining principal
	remainingPrincipal := loan.PrincipalAmount - loan.TotalInvestedAmount
	if investment.InvestmentAmount > remainingPrincipal {
		return fmt.Errorf("investment amount exceeds remaining principal. Remaining: %f", remainingPrincipal)
	}
	
	// Check if investor already invested in this loan
	existingInvestment, err := s.loanInvestmentRepo.GetByLoanAndInvestor(ctx, loanID, investment.InvestorID)
	if err == nil && existingInvestment != nil {
		return errors.New("investor already invested in this loan")
	}
	
	// Create investment record
	investment.LoanID = loanID
	err = s.loanInvestmentRepo.Create(ctx, investment)
	if err != nil {
		return fmt.Errorf("failed to create investment: %w", err)
	}
	
	// Update total invested amount in loan
	newTotal := loan.TotalInvestedAmount + investment.InvestmentAmount
	err = s.loanRepo.UpdateTotalInvestedAmount(ctx, loanID, newTotal)
	if err != nil {
		return fmt.Errorf("failed to update total invested amount: %w", err)
	}
	
	// Check if loan is fully invested
	if newTotal >= loan.PrincipalAmount {
		// Update loan state to invested
		err = s.loanRepo.UpdateState(ctx, loanID, "invested")
		if err != nil {
			return fmt.Errorf("failed to update loan state: %w", err)
		}
		
		// Add state transition to history
		stateHistory := &models.LoanStateHistory{
			LoanID:           loanID,
			PreviousState:    loan.CurrentState,
			NewState:         "invested",
			TransitionReason: "Loan fully invested",
		}
		
		err = s.loanStateHistoryRepo.Create(ctx, stateHistory)
		if err != nil {
			return fmt.Errorf("failed to create state history: %w", err)
		}
		
		// Send investment confirmation emails to all investors
		investments, err := s.loanInvestmentRepo.GetByLoanID(ctx, loanID)
		if err != nil {
			return fmt.Errorf("failed to get loan investments: %w", err)
		}
		
		for _, inv := range investments {
			investor, err := s.investorRepo.GetByID(ctx, inv.InvestorID)
			if err != nil {
				continue // Log error but continue with other investors
			}
			
			// Send investment confirmation email
			err = s.emailService.SendInvestmentConfirmation(ctx, investor.Email, loan.AgreementLetterLink, fmt.Sprintf("Loan %s has been fully invested", loan.LoanID))
			if err != nil {
				// Log error but continue with other investors
			}
		}
	}
	
	return nil
}

func (s *loanServiceImpl) DisburseLoan(ctx context.Context, loanID int, disbursementData *models.LoanDisbursement) error {
	// Get the loan
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}
	
	// Check if loan is in invested state
	if loan.CurrentState != "invested" {
		return errors.New("loan must be in invested state to be disbursed")
	}
	
	// Check if total invested amount equals principal amount
	if loan.TotalInvestedAmount != loan.PrincipalAmount {
		return errors.New("total invested amount must equal principal amount for disbursement")
	}
	
	// Validate disbursement data
	if disbursementData.FieldOfficerEmployeeID == "" {
		return errors.New("field officer employee ID is required")
	}
	
	if disbursementData.AgreementLetterSignedUrl == "" {
		return errors.New("signed agreement letter URL is required")
	}
	
	// Create loan disbursement record
	disbursementData.LoanID = loanID
	err = s.loanDisbursementRepo.Create(ctx, disbursementData)
	if err != nil {
		return fmt.Errorf("failed to create loan disbursement: %w", err)
	}
	
	// Update loan state to disbursed
	err = s.loanRepo.UpdateState(ctx, loanID, "disbursed")
	if err != nil {
		return fmt.Errorf("failed to update loan state: %w", err)
	}
	
	// Add state transition to history
	stateHistory := &models.LoanStateHistory{
		LoanID:           loanID,
		PreviousState:    loan.CurrentState,
		NewState:         "disbursed",
		TransitionReason: "Loan disbursed to borrower",
	}
	
	err = s.loanStateHistoryRepo.Create(ctx, stateHistory)
	if err != nil {
		return fmt.Errorf("failed to create state history: %w", err)
	}
	
	return nil
}

func (s *loanServiceImpl) GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error) {
	return s.loanRepo.GetTotalInvestedAmount(ctx, loanID)
}

func (s *loanServiceImpl) CanTransitionToState(ctx context.Context, loanID int, newState string) (bool, error) {
	loan, err := s.loanRepo.GetByID(ctx, loanID)
	if err != nil {
		return false, err
	}
	
	currentState := loan.CurrentState
	
	// Define valid state transitions
	validTransitions := map[string][]string{
		"proposed": {"approved"},
		"approved": {"invested"},
		"invested": {"disbursed"},
		"disbursed": {}, // No further transitions allowed
	}
	
	validStates, exists := validTransitions[currentState]
	if !exists {
		return false, fmt.Errorf("invalid current state: %s", currentState)
	}
	
	for _, state := range validStates {
		if state == newState {
			return true, nil
		}
	}
	
	return false, nil
}