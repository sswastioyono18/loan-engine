package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/sswastioyono18/loan-engine/internal/models"
	"github.com/sswastioyono18/loan-engine/internal/repositories/mocks"
	mocks2 "github.com/sswastioyono18/loan-engine/pkg/external/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateLoan(t *testing.T) {
	mockLoanRepo := mocks.NewLoanRepository(t)
	mockApprovalRepo := mocks.NewLoanApprovalRepository(t)
	mockDisbursementRepo := mocks.NewLoanDisbursementRepository(t)
	mockInvestmentRepo := mocks.NewLoanInvestmentRepository(t)
	mockStateHistoryRepo := mocks.NewLoanStateHistoryRepository(t)
	mockInvestorRepo := mocks.NewInvestorRepository(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loan := &models.Loan{
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
	}

	mockLoanRepo.On("Create", context.Background(), loan).Return(nil)

	err := service.CreateLoan(context.Background(), loan)

	assert.NoError(t, err)
	assert.Equal(t, "proposed", loan.CurrentState)
}

func TestApproveLoan(t *testing.T) {
	mockLoanRepo := mocks.NewLoanRepository(t)
	mockApprovalRepo := mocks.NewLoanApprovalRepository(t)
	mockDisbursementRepo := mocks.NewLoanDisbursementRepository(t)
	mockInvestmentRepo := mocks.NewLoanInvestmentRepository(t)
	mockStateHistoryRepo := mocks.NewLoanStateHistoryRepository(t)
	mockInvestorRepo := mocks.NewInvestorRepository(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "proposed",
	}

	approval := &models.LoanApproval{
		FieldValidatorEmployeeID: "emp001",
		ProofImageUrl:            "https://example.com/proof.jpg",
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)
	mockApprovalRepo.On("Create", context.Background(), approval).Return(nil)
	mockLoanRepo.On("UpdateState", context.Background(), loanID, "approved").Return(nil)
	mockStateHistoryRepo.On("Create", context.Background(), mock.Anything).Return(nil)

	err := service.ApproveLoan(context.Background(), loanID, approval)

	assert.NoError(t, err)
}

func TestApproveLoanInvalidState(t *testing.T) {
	mockLoanRepo := mocks.NewLoanRepository(t)
	mockApprovalRepo := mocks.NewLoanApprovalRepository(t)
	mockDisbursementRepo := mocks.NewLoanDisbursementRepository(t)
	mockInvestmentRepo := mocks.NewLoanInvestmentRepository(t)
	mockStateHistoryRepo := mocks.NewLoanStateHistoryRepository(t)
	mockInvestorRepo := mocks.NewInvestorRepository(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "approved", // Already approved
	}

	approval := &models.LoanApproval{
		FieldValidatorEmployeeID: "emp001",
		ProofImageUrl:            "https://example.com/proof.jpg",
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)

	err := service.ApproveLoan(context.Background(), loanID, approval)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loan must be in proposed state to be approved")
}

func TestInvestInLoan(t *testing.T) {
	mockLoanRepo := mocks.NewLoanRepository(t)
	mockApprovalRepo := mocks.NewLoanApprovalRepository(t)
	mockDisbursementRepo := mocks.NewLoanDisbursementRepository(t)
	mockInvestmentRepo := mocks.NewLoanInvestmentRepository(t)
	mockStateHistoryRepo := mocks.NewLoanStateHistoryRepository(t)
	mockInvestorRepo := mocks.NewInvestorRepository(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "approved",
		TotalInvestedAmount: 0.0,
	}

	investment := &models.LoanInvestment{
		InvestorID:       1,
		InvestmentAmount: 5000.0,
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)
	mockInvestmentRepo.On("GetByLoanAndInvestor", context.Background(), loanID, 1).Return(nil, errors.New("not found"))
	mockInvestmentRepo.On("Create", context.Background(), investment).Return(nil)
	mockLoanRepo.On("UpdateTotalInvestedAmount", context.Background(), loanID, 5000.0).Return(nil)

	err := service.InvestInLoan(context.Background(), loanID, investment)

	assert.NoError(t, err)
}

func TestInvestInLoanExceedsPrincipal(t *testing.T) {
	mockLoanRepo := mocks.NewLoanRepository(t)
	mockApprovalRepo := mocks.NewLoanApprovalRepository(t)
	mockDisbursementRepo := mocks.NewLoanDisbursementRepository(t)
	mockInvestmentRepo := mocks.NewLoanInvestmentRepository(t)
	mockStateHistoryRepo := mocks.NewLoanStateHistoryRepository(t)
	mockInvestorRepo := mocks.NewInvestorRepository(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "approved",
		TotalInvestedAmount: 5000.0,
	}

	investment := &models.LoanInvestment{
		InvestorID:       1,
		InvestmentAmount: 6000.0, // Exceeds remaining principal (10000 - 5000 = 5000)
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)

	err := service.InvestInLoan(context.Background(), loanID, investment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "investment amount exceeds remaining principal")
}

func TestDisburseLoan(t *testing.T) {
	mockLoanRepo := mocks.NewLoanRepository(t)
	mockApprovalRepo := mocks.NewLoanApprovalRepository(t)
	mockDisbursementRepo := mocks.NewLoanDisbursementRepository(t)
	mockInvestmentRepo := mocks.NewLoanInvestmentRepository(t)
	mockStateHistoryRepo := mocks.NewLoanStateHistoryRepository(t)
	mockInvestorRepo := mocks.NewInvestorRepository(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "invested",
		TotalInvestedAmount: 10000.0,
	}

	disbursement := &models.LoanDisbursement{
		FieldOfficerEmployeeID:   "emp002",
		AgreementLetterSignedUrl: "https://example.com/signed-agreement.pdf",
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)
	mockDisbursementRepo.On("Create", context.Background(), disbursement).Return(nil)
	mockLoanRepo.On("UpdateState", context.Background(), loanID, "disbursed").Return(nil)
	mockStateHistoryRepo.On("Create", context.Background(), mock.Anything).Return(nil)

	err := service.DisburseLoan(context.Background(), loanID, disbursement)

	assert.NoError(t, err)
}

func TestDisburseLoanInvalidState(t *testing.T) {
	mockLoanRepo := mocks.NewLoanRepository(t)
	mockApprovalRepo := mocks.NewLoanApprovalRepository(t)
	mockDisbursementRepo := mocks.NewLoanDisbursementRepository(t)
	mockInvestmentRepo := mocks.NewLoanInvestmentRepository(t)
	mockStateHistoryRepo := mocks.NewLoanStateHistoryRepository(t)
	mockInvestorRepo := mocks.NewInvestorRepository(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "proposed", // Not invested yet
		TotalInvestedAmount: 0.0,
	}

	disbursement := &models.LoanDisbursement{
		FieldOfficerEmployeeID:   "emp002",
		AgreementLetterSignedUrl: "https://example.com/signed-agreement.pdf",
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)

	err := service.DisburseLoan(context.Background(), loanID, disbursement)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loan must be in invested state to be disbursed")
}
