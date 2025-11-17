package services

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/kitabisa/loan-engine/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockLoanRepository struct {
	mock.Mock
}

func (m *MockLoanRepository) Create(ctx context.Context, loan *models.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GetByID(ctx context.Context, id int) (*models.Loan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) Update(ctx context.Context, loan *models.Loan) error {
	args := m.Called(ctx, loan)
	return args.Error(0)
}

func (m *MockLoanRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLoanRepository) List(ctx context.Context, state *string, offset, limit int) ([]*models.Loan, error) {
	args := m.Called(ctx, state, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetByState(ctx context.Context, state string) ([]*models.Loan, error) {
	args := m.Called(ctx, state)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetTotalInvestmentAmount(ctx context.Context, loanID int) (float64, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockLoanRepository) GetByLoanID(ctx context.Context, loanID string) (*models.Loan, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepository) UpdateState(ctx context.Context, id int, newState string) error {
	args := m.Called(ctx, id, newState)
	return args.Error(0)
}

func (m *MockLoanRepository) UpdateTotalInvestedAmount(ctx context.Context, loanID int, amount float64) error {
	args := m.Called(ctx, loanID, amount)
	return args.Error(0)
}

func (m *MockLoanRepository) GetTotalInvestedAmount(ctx context.Context, loanID int) (float64, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(float64), args.Error(1)
}

type MockLoanApprovalRepository struct {
	mock.Mock
}

func (m *MockLoanApprovalRepository) Create(ctx context.Context, approval *models.LoanApproval) error {
	args := m.Called(ctx, approval)
	return args.Error(0)
}

func (m *MockLoanApprovalRepository) GetByLoanID(ctx context.Context, loanID int) (*models.LoanApproval, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanApproval), args.Error(1)
}

func (m *MockLoanApprovalRepository) GetByID(ctx context.Context, id int) (*models.LoanApproval, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanApproval), args.Error(1)
}

func (m *MockLoanApprovalRepository) Update(ctx context.Context, approval *models.LoanApproval) error {
	args := m.Called(ctx, approval)
	return args.Error(0)
}

func (m *MockLoanApprovalRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockLoanDisbursementRepository struct {
	mock.Mock
}

func (m *MockLoanDisbursementRepository) Create(ctx context.Context, disbursement *models.LoanDisbursement) error {
	args := m.Called(ctx, disbursement)
	return args.Error(0)
}

func (m *MockLoanDisbursementRepository) GetByLoanID(ctx context.Context, loanID int) (*models.LoanDisbursement, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanDisbursement), args.Error(1)
}

func (m *MockLoanDisbursementRepository) GetByID(ctx context.Context, id int) (*models.LoanDisbursement, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanDisbursement), args.Error(1)
}

func (m *MockLoanDisbursementRepository) Update(ctx context.Context, disbursement *models.LoanDisbursement) error {
	args := m.Called(ctx, disbursement)
	return args.Error(0)
}

func (m *MockLoanDisbursementRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockLoanInvestmentRepository struct {
	mock.Mock
}

func (m *MockLoanInvestmentRepository) Create(ctx context.Context, investment *models.LoanInvestment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockLoanInvestmentRepository) GetByID(ctx context.Context, id int) (*models.LoanInvestment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanInvestment), args.Error(1)
}

func (m *MockLoanInvestmentRepository) GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanInvestment, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.LoanInvestment), args.Error(1)
}

func (m *MockLoanInvestmentRepository) GetByInvestorID(ctx context.Context, investorID int) ([]*models.LoanInvestment, error) {
	args := m.Called(ctx, investorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.LoanInvestment), args.Error(1)
}

func (m *MockLoanInvestmentRepository) GetByLoanAndInvestor(ctx context.Context, loanID, investorID int) (*models.LoanInvestment, error) {
	args := m.Called(ctx, loanID, investorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanInvestment), args.Error(1)
}

func (m *MockLoanInvestmentRepository) Update(ctx context.Context, investment *models.LoanInvestment) error {
	args := m.Called(ctx, investment)
	return args.Error(0)
}

func (m *MockLoanInvestmentRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockLoanInvestmentRepository) GetTotalInvestedAmountByLoan(ctx context.Context, loanID int) (float64, error) {
	args := m.Called(ctx, loanID)
	return args.Get(0).(float64), args.Error(1)
}

type MockLoanStateHistoryRepository struct {
	mock.Mock
}

func (m *MockLoanStateHistoryRepository) Create(ctx context.Context, history *models.LoanStateHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

func (m *MockLoanStateHistoryRepository) GetByLoanID(ctx context.Context, loanID int) ([]*models.LoanStateHistory, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.LoanStateHistory), args.Error(1)
}

func (m *MockLoanStateHistoryRepository) GetLatestByLoanID(ctx context.Context, loanID int) (*models.LoanStateHistory, error) {
	args := m.Called(ctx, loanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.LoanStateHistory), args.Error(1)
}

func (m *MockLoanStateHistoryRepository) List(ctx context.Context, loanID int, offset, limit int) ([]*models.LoanStateHistory, error) {
	args := m.Called(ctx, loanID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.LoanStateHistory), args.Error(1)
}

type MockInvestorRepository struct {
	mock.Mock
}

func (m *MockInvestorRepository) GetByID(ctx context.Context, id int) (*models.Investor, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Investor), args.Error(1)
}

func (m *MockInvestorRepository) Create(ctx context.Context, investor *models.Investor) error {
	args := m.Called(ctx, investor)
	return args.Error(0)
}

func (m *MockInvestorRepository) GetByInvestorID(ctx context.Context, investorID string) (*models.Investor, error) {
	args := m.Called(ctx, investorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Investor), args.Error(1)
}

func (m *MockInvestorRepository) GetByEmail(ctx context.Context, email string) (*models.Investor, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Investor), args.Error(1)
}

func (m *MockInvestorRepository) Update(ctx context.Context, investor *models.Investor) error {
	args := m.Called(ctx, investor)
	return args.Error(0)
}

func (m *MockInvestorRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockInvestorRepository) List(ctx context.Context, offset, limit int) ([]*models.Investor, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Investor), args.Error(1)
}

// Mock external services
type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func (m *MockEmailService) SendInvestmentConfirmation(ctx context.Context, to, agreementLink, message string) error {
	args := m.Called(ctx, to, agreementLink, message)
	return args.Error(0)
}

func (m *MockEmailService) SendDisbursementNotification(ctx context.Context, to, loanDetails string) error {
	args := m.Called(ctx, to, loanDetails)
	return args.Error(0)
}

func (m *MockEmailService) SendApprovalNotification(ctx context.Context, to, loanDetails string) error {
	args := m.Called(ctx, to, loanDetails)
	return args.Error(0)
}

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) UploadFile(ctx context.Context, file io.Reader, fileName, contentType string) (string, error) {
	args := m.Called(ctx, file, fileName, contentType)
	return args.String(0), args.Error(1)
}

func (m *MockStorageService) DownloadFile(ctx context.Context, fileID string) (io.ReadCloser, error) {
	args := m.Called(ctx, fileID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) DeleteFile(ctx context.Context, fileID string) error {
	args := m.Called(ctx, fileID)
	return args.Error(0)
}

func (m *MockStorageService) GetFileURL(ctx context.Context, fileID string) (string, error) {
	args := m.Called(ctx, fileID)
	return args.String(0), args.Error(1)
}

func TestCreateLoan(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockLoanApprovalRepository)
	mockDisbursementRepo := new(MockLoanDisbursementRepository)
	mockInvestmentRepo := new(MockLoanInvestmentRepository)
	mockStateHistoryRepo := new(MockLoanStateHistoryRepository)
	mockInvestorRepo := new(MockInvestorRepository)
	mockEmailService := new(MockEmailService)
	mockStorageService := new(MockStorageService)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loan := &models.Loan{
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: "https://example.com/agreement.pdf",
	}

	mockLoanRepo.On("Create", context.Background(), loan).Return(nil)

	err := service.CreateLoan(context.Background(), loan)

	assert.NoError(t, err)
	assert.Equal(t, "proposed", loan.CurrentState)
	mockLoanRepo.AssertExpectations(t)
}

func TestApproveLoan(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockLoanApprovalRepository)
	mockDisbursementRepo := new(MockLoanDisbursementRepository)
	mockInvestmentRepo := new(MockLoanInvestmentRepository)
	mockStateHistoryRepo := new(MockLoanStateHistoryRepository)
	mockInvestorRepo := new(MockInvestorRepository)
	mockEmailService := new(MockEmailService)
	mockStorageService := new(MockStorageService)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: "https://example.com/agreement.pdf",
		CurrentState:        "proposed",
	}

	approval := &models.LoanApproval{
		FieldValidatorEmployeeID: "emp001",
		ProofImageUrl:            "https://example.com/proof.jpg",
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)
	mockApprovalRepo.On("Create", context.Background(), approval).Return(nil)
	mockLoanRepo.On("UpdateState", context.Background(), loanID, "approved").Return(nil)
	mockStateHistoryRepo.On("Create", context.Background(), mock.AnythingOfType("*models.LoanStateHistory")).Return(nil)

	err := service.ApproveLoan(context.Background(), loanID, approval)

	assert.NoError(t, err)
	mockLoanRepo.AssertExpectations(t)
	mockApprovalRepo.AssertExpectations(t)
	mockStateHistoryRepo.AssertExpectations(t)
}

func TestApproveLoanInvalidState(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockLoanApprovalRepository)
	mockDisbursementRepo := new(MockLoanDisbursementRepository)
	mockInvestmentRepo := new(MockLoanInvestmentRepository)
	mockStateHistoryRepo := new(MockLoanStateHistoryRepository)
	mockInvestorRepo := new(MockInvestorRepository)
	mockEmailService := new(MockEmailService)
	mockStorageService := new(MockStorageService)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: "https://example.com/agreement.pdf",
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
	mockLoanRepo.AssertExpectations(t)
}

func TestInvestInLoan(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockLoanApprovalRepository)
	mockDisbursementRepo := new(MockLoanDisbursementRepository)
	mockInvestmentRepo := new(MockLoanInvestmentRepository)
	mockStateHistoryRepo := new(MockLoanStateHistoryRepository)
	mockInvestorRepo := new(MockInvestorRepository)
	mockEmailService := new(MockEmailService)
	mockStorageService := new(MockStorageService)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: "https://example.com/agreement.pdf",
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
	mockLoanRepo.AssertExpectations(t)
	mockInvestmentRepo.AssertExpectations(t)
}

func TestInvestInLoanExceedsPrincipal(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockLoanApprovalRepository)
	mockDisbursementRepo := new(MockLoanDisbursementRepository)
	mockInvestmentRepo := new(MockLoanInvestmentRepository)
	mockStateHistoryRepo := new(MockLoanStateHistoryRepository)
	mockInvestorRepo := new(MockInvestorRepository)
	mockEmailService := new(MockEmailService)
	mockStorageService := new(MockStorageService)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: "https://example.com/agreement.pdf",
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
	mockLoanRepo.AssertExpectations(t)
}

func TestDisburseLoan(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockLoanApprovalRepository)
	mockDisbursementRepo := new(MockLoanDisbursementRepository)
	mockInvestmentRepo := new(MockLoanInvestmentRepository)
	mockStateHistoryRepo := new(MockLoanStateHistoryRepository)
	mockInvestorRepo := new(MockInvestorRepository)
	mockEmailService := new(MockEmailService)
	mockStorageService := new(MockStorageService)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: "https://example.com/agreement.pdf",
		CurrentState:        "invested",
		TotalInvestedAmount: 10000.0,
	}

	disbursement := &models.LoanDisbursement{
		FieldOfficerEmployeeID:      "emp002",
		AgreementLetterSignedUrl:    "https://example.com/signed-agreement.pdf",
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)
	mockDisbursementRepo.On("Create", context.Background(), disbursement).Return(nil)
	mockLoanRepo.On("UpdateState", context.Background(), loanID, "disbursed").Return(nil)
	mockStateHistoryRepo.On("Create", context.Background(), mock.AnythingOfType("*models.LoanStateHistory")).Return(nil)

	err := service.DisburseLoan(context.Background(), loanID, disbursement)

	assert.NoError(t, err)
	mockLoanRepo.AssertExpectations(t)
	mockDisbursementRepo.AssertExpectations(t)
	mockStateHistoryRepo.AssertExpectations(t)
}

func TestDisburseLoanInvalidState(t *testing.T) {
	mockLoanRepo := new(MockLoanRepository)
	mockApprovalRepo := new(MockLoanApprovalRepository)
	mockDisbursementRepo := new(MockLoanDisbursementRepository)
	mockInvestmentRepo := new(MockLoanInvestmentRepository)
	mockStateHistoryRepo := new(MockLoanStateHistoryRepository)
	mockInvestorRepo := new(MockInvestorRepository)
	mockEmailService := new(MockEmailService)
	mockStorageService := new(MockStorageService)

	service := NewLoanService(mockLoanRepo, mockApprovalRepo, mockDisbursementRepo, mockInvestmentRepo, mockStateHistoryRepo, mockInvestorRepo, mockEmailService, mockStorageService)

	loanID := 1
	loan := &models.Loan{
		ID:                  loanID,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: "https://example.com/agreement.pdf",
		CurrentState:        "proposed", // Not invested yet
		TotalInvestedAmount: 0.0,
	}

	disbursement := &models.LoanDisbursement{
		FieldOfficerEmployeeID:      "emp002",
		AgreementLetterSignedUrl:    "https://example.com/signed-agreement.pdf",
	}

	mockLoanRepo.On("GetByID", context.Background(), loanID).Return(loan, nil)

	err := service.DisburseLoan(context.Background(), loanID, disbursement)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loan must be in invested state to be disbursed")
	mockLoanRepo.AssertExpectations(t)
}