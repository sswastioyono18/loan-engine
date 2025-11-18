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

func TestInvestInLoanSendsEmailNotifications(t *testing.T) {
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
		TotalInvestedAmount: 5000.0, // Need 5000 more to reach full amount
		LoanID:              "LOAN001",
	}

	investment := &models.LoanInvestment{
		InvestorID:       1,
		InvestmentAmount: 5000.0, // This will make total invested = 10000 (equal to principal)
	}

	investor := &models.Investor{
		ID:       1,
		Email:    "investor@example.com",
		FullName: "Test Investor",
	}

	loanInvestments := []*models.LoanInvestment{
		{
			ID:             1,
			InvestorID:     1,
			InvestmentAmount: 5000.0,
		},
	}

	// Set up mocks
	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)
	mockInvestmentRepo.On("GetByLoanAndInvestor", context.Background(), loanID, 1).Return(nil, errors.New("not found"))
	mockInvestmentRepo.On("Create", context.Background(), investment).Return(nil)
	mockLoanRepo.On("UpdateTotalInvestedAmount", context.Background(), loanID, 10000.0).Return(nil)
	mockLoanRepo.On("UpdateState", context.Background(), loanID, "invested").Return(nil)
	mockStateHistoryRepo.On("Create", context.Background(), mock.Anything).Return(nil)
	mockInvestmentRepo.On("GetByLoanID", context.Background(), loanID).Return(loanInvestments, nil)
	mockInvestorRepo.On("GetByID", context.Background(), 1).Return(investor, nil)
	mockEmailService.On("SendInvestmentConfirmation", context.Background(), "investor@example.com", "https://example.com/agreement.pdf", "Loan LOAN001 has been fully invested").Return(nil)

	err := service.InvestInLoan(context.Background(), loanID, investment)

	assert.NoError(t, err)
	// Note: The loan object in the test won't be updated by the service method, so we can't check loan.CurrentState directly
	// The state update is handled by the repository, which is mocked
	mockEmailService.AssertExpectations(t)
}

func TestCanTransitionToState(t *testing.T) {
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

	// Test valid transitions
	tests := []struct {
		name           string
		currentState   string
		targetState    string
		expectedResult bool
		shouldError    bool
	}{
		{
			name:           "proposed to approved - valid",
			currentState:   "proposed",
			targetState:    "approved",
			expectedResult: true,
			shouldError:    false,
		},
		{
			name:           "approved to invested - valid",
			currentState:   "approved",
			targetState:    "invested",
			expectedResult: true,
			shouldError:    false,
		},
		{
			name:           "invested to disbursed - valid",
			currentState:   "invested",
			targetState:    "disbursed",
			expectedResult: true,
			shouldError:    false,
		},
		{
			name:           "proposed to invested - invalid",
			currentState:   "proposed",
			targetState:    "invested",
			expectedResult: false,
			shouldError:    false,
		},
		{
			name:           "approved to disbursed - invalid",
			currentState:   "approved",
			targetState:    "disbursed",
			expectedResult: false,
			shouldError:    false,
		},
		{
			name:           "disbursed to any - invalid",
			currentState:   "disbursed",
			targetState:    "approved",
			expectedResult: false,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loan := &models.Loan{
				ID:           loanID,
				CurrentState: tt.currentState,
			}

			mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil).Once()

			canTransition, err := service.CanTransitionToState(context.Background(), loanID, tt.targetState)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, canTransition)
			}
		})
	}
}

func TestStateHistoryRecordedDuringTransitions(t *testing.T) {
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

	// Test state history during approval
	t.Run("approval state history", func(t *testing.T) {
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

		mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil).Once()
		mockApprovalRepo.On("Create", context.Background(), approval).Return(nil).Once()
		mockLoanRepo.On("UpdateState", context.Background(), loanID, "approved").Return(nil).Once()
		mockStateHistoryRepo.On("Create", context.Background(), mock.MatchedBy(func(history *models.LoanStateHistory) bool {
			return history.LoanID == loanID &&
				   history.PreviousState == "proposed" &&
				   history.NewState == "approved" &&
				   history.TransitionReason == "Loan approved by staff"
		})).Return(nil).Once()

		err := service.ApproveLoan(context.Background(), loanID, approval)

		assert.NoError(t, err)
		mockStateHistoryRepo.AssertExpectations(t)
	})

	// Reset mocks for next test
	mockStateHistoryRepo.ExpectedCalls = nil
	mockLoanRepo.ExpectedCalls = nil
	mockApprovalRepo.ExpectedCalls = nil

	// Test state history during investment to "invested" state
	t.Run("investment state history", func(t *testing.T) {
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
			InvestmentAmount: 5000.0, // This will make total invested = 10000 (equal to principal)
		}

		mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil).Once()
		mockInvestmentRepo.On("GetByLoanAndInvestor", context.Background(), loanID, 1).Return(nil, errors.New("not found")).Once()
		mockInvestmentRepo.On("Create", context.Background(), investment).Return(nil).Once()
		mockLoanRepo.On("UpdateTotalInvestedAmount", context.Background(), loanID, 10000.0).Return(nil).Once()
		mockLoanRepo.On("UpdateState", context.Background(), loanID, "invested").Return(nil).Once()
		mockStateHistoryRepo.On("Create", context.Background(), mock.MatchedBy(func(history *models.LoanStateHistory) bool {
			return history.LoanID == loanID &&
				   history.PreviousState == "approved" &&
				   history.NewState == "invested" &&
				   history.TransitionReason == "Loan fully invested"
		})).Return(nil).Once()
		mockInvestmentRepo.On("GetByLoanID", context.Background(), loanID).Return([]*models.LoanInvestment{investment}, nil).Once()
		mockInvestorRepo.On("GetByID", context.Background(), 1).Return(&models.Investor{ID: 1, Email: "investor@example.com"}, nil).Once()
		mockEmailService.On("SendInvestmentConfirmation", context.Background(), "investor@example.com", "https://example.com/agreement.pdf", mock.Anything).Return(nil).Once()

		err := service.InvestInLoan(context.Background(), loanID, investment)

		assert.NoError(t, err)
		// Note: The loan object in the test won't be updated by the service method, so we can't check loan.CurrentState directly
		// The state update is handled by the repository, which is mocked
		mockStateHistoryRepo.AssertExpectations(t)
	})

	// Reset mocks for next test
	mockStateHistoryRepo.ExpectedCalls = nil
	mockLoanRepo.ExpectedCalls = nil
	mockInvestmentRepo.ExpectedCalls = nil
	mockEmailService.ExpectedCalls = nil
	mockInvestorRepo.ExpectedCalls = nil

	// Test state history during disbursement
	t.Run("disbursement state history", func(t *testing.T) {
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

		mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil).Once()
		mockDisbursementRepo.On("Create", context.Background(), disbursement).Return(nil).Once()
		mockLoanRepo.On("UpdateState", context.Background(), loanID, "disbursed").Return(nil).Once()
		mockStateHistoryRepo.On("Create", context.Background(), mock.MatchedBy(func(history *models.LoanStateHistory) bool {
			return history.LoanID == loanID &&
				   history.PreviousState == "invested" &&
				   history.NewState == "disbursed" &&
				   history.TransitionReason == "Loan disbursed to borrower"
		})).Return(nil).Once()

		err := service.DisburseLoan(context.Background(), loanID, disbursement)

		assert.NoError(t, err)
		mockStateHistoryRepo.AssertExpectations(t)
	})
}

func TestMultipleInvestorsInSameLoan(t *testing.T) {
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
		LoanID:              "LOAN001",
	}

	investor1 := &models.Investor{
		ID:       1,
		Email:    "investor1@example.com",
		FullName: "Investor 1",
	}

	investor2 := &models.Investor{
		ID:       2,
		Email:    "investor2@example.com",
		FullName: "Investor 2",
	}

	investment1 := &models.LoanInvestment{
		InvestorID:       1,
		InvestmentAmount: 6000.0, // First investment
	}

	investment2 := &models.LoanInvestment{
		InvestorID:       2,
		InvestmentAmount: 4000.0, // Second investment to reach principal
	}

	loanInvestments := []*models.LoanInvestment{
		{
			ID:             1,
			InvestorID:     1,
			InvestmentAmount: 6000.0,
		},
		{
			ID:             2,
			InvestorID:     2,
			InvestmentAmount: 4000.0,
		},
	}

	// First investment - should succeed
	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil).Once()
	mockInvestmentRepo.On("GetByLoanAndInvestor", context.Background(), loanID, 1).Return(nil, errors.New("not found")).Once()
	mockInvestmentRepo.On("Create", context.Background(), investment1).Return(nil).Once()
	mockLoanRepo.On("UpdateTotalInvestedAmount", context.Background(), loanID, 6000.0).Return(nil).Once()

	err := service.InvestInLoan(context.Background(), loanID, investment1)
	assert.NoError(t, err)

	// Second investment - should make loan fully invested and trigger emails
	// Create a new loan object for the second call to GetByID with updated TotalInvestedAmount
	loanAfterFirstInvestment := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "approved",
		TotalInvestedAmount: 6000.0, // Updated after first investment
		LoanID:              "LOAN001",
	}
	
	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loanAfterFirstInvestment, nil).Once()
	mockInvestmentRepo.On("GetByLoanAndInvestor", context.Background(), loanID, 2).Return(nil, errors.New("not found")).Once()
	mockInvestmentRepo.On("Create", context.Background(), investment2).Return(nil).Once()
	mockLoanRepo.On("UpdateTotalInvestedAmount", context.Background(), loanID, 10000.0).Return(nil).Once()
	mockLoanRepo.On("UpdateState", context.Background(), loanID, "invested").Return(nil).Once()
	mockStateHistoryRepo.On("Create", context.Background(), mock.MatchedBy(func(history *models.LoanStateHistory) bool {
		return history.LoanID == loanID &&
			   history.PreviousState == "approved" &&
			   history.NewState == "invested" &&
			   history.TransitionReason == "Loan fully invested"
	})).Return(nil).Once()
	mockInvestmentRepo.On("GetByLoanID", context.Background(), loanID).Return(loanInvestments, nil).Once()
	mockInvestorRepo.On("GetByID", context.Background(), 1).Return(investor1, nil).Once()
	mockInvestorRepo.On("GetByID", context.Background(), 2).Return(investor2, nil).Once()
	mockEmailService.On("SendInvestmentConfirmation", context.Background(), "investor1@example.com", "https://example.com/agreement.pdf", "Loan LOAN001 has been fully invested").Return(nil).Once()
	mockEmailService.On("SendInvestmentConfirmation", context.Background(), "investor2@example.com", "https://example.com/agreement.pdf", "Loan LOAN001 has been fully invested").Return(nil).Once()

	err = service.InvestInLoan(context.Background(), loanID, investment2)

	assert.NoError(t, err)
	// Note: The loan object in the test won't be updated by the service method, so we can't check loan.CurrentState directly
	// The state update is handled by the repository, which is mocked
	mockEmailService.AssertExpectations(t)
}
