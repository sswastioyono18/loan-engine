package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sswastioyono18/loan-engine/internal/models"
	"github.com/sswastioyono18/loan-engine/internal/services/mocks"
	mocks2 "github.com/sswastioyono18/loan-engine/pkg/external/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoanHandlerCreateLoan(t *testing.T) {
	mockLoanService := mocks.NewLoanService(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	handler := NewLoanHandler(mockLoanService, mockEmailService, mockStorageService)

	loanReq := map[string]interface{}{
		"borrower_id":           1,
		"principal_amount":      10000.0,
		"rate":                  0.05,
		"roi":                   0.08,
		"agreement_letter_link": "https://example.com/agreement.pdf",
	}

	loanReqBytes, _ := json.Marshal(loanReq)

	req, _ := http.NewRequest("POST", "/api/v1/loans", bytes.NewBuffer(loanReqBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Create expected model that the service will receive
	expectedModel := &models.Loan{
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
	}

	mockLoanService.On("CreateLoan", mock.Anything, mock.MatchedBy(func(loan *models.Loan) bool {
		return loan.BorrowerID == expectedModel.BorrowerID &&
			loan.PrincipalAmount == expectedModel.PrincipalAmount &&
			loan.Rate == expectedModel.Rate &&
			loan.ROI == expectedModel.ROI &&
			loan.AgreementLetterLink.String == expectedModel.AgreementLetterLink.String
	})).Return(nil)

	handler.CreateLoan(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code) // Should be 200, not 201
	mockLoanService.AssertExpectations(t)
}

func TestLoanHandlerGetLoanByID(t *testing.T) {
	mockLoanService := mocks.NewLoanService(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	handler := NewLoanHandler(mockLoanService, mockEmailService, mockStorageService)

	loan := &models.Loan{
		ID:                  1,
		BorrowerID:          1,
		PrincipalAmount:     10000.0,
		Rate:                0.05,
		ROI:                 0.08,
		AgreementLetterLink: sql.NullString{String: "https://example.com/agreement.pdf", Valid: true},
		CurrentState:        "proposed",
	}

	req, _ := http.NewRequest("GET", "/api/v1/loans/1", nil)
	rr := httptest.NewRecorder()

	// Set up chi URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockLoanService.On("GetLoanByID", mock.Anything, 1).Return(loan, nil)

	handler.GetLoanByID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)
	
	// The response is wrapped in a Response struct with a Data field
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected response.data to be a map, got %T", response["data"])
	}
	
	assert.Equal(t, float64(1), data["id"])
	mockLoanService.AssertExpectations(t)
}

func TestLoanHandlerApproveLoan(t *testing.T) {
	mockLoanService := mocks.NewLoanService(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	handler := NewLoanHandler(mockLoanService, mockEmailService, mockStorageService)

	approvalReq := &models.LoanApproval{
		FieldValidatorEmployeeID: "emp001",
		ProofImageUrl:            "https://example.com/proof.jpg",
	}

	approvalReqBytes, _ := json.Marshal(approvalReq)

	req, _ := http.NewRequest("POST", "/api/v1/loans/1/approve", bytes.NewBuffer(approvalReqBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Set up chi URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockLoanService.On("ApproveLoan", mock.Anything, 1, mock.AnythingOfType("*models.LoanApproval")).Return(nil)

	handler.ApproveLoan(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockLoanService.AssertExpectations(t)
}

func TestLoanHandlerInvestInLoan(t *testing.T) {
	mockLoanService := mocks.NewLoanService(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	handler := NewLoanHandler(mockLoanService, mockEmailService, mockStorageService)

	investmentReq := &models.LoanInvestment{
		InvestorID:       1,
		InvestmentAmount: 5000.0,
	}

	investmentReqBytes, _ := json.Marshal(investmentReq)

	req, _ := http.NewRequest("POST", "/api/v1/loans/1/invest", bytes.NewBuffer(investmentReqBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Set up chi URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockLoanService.On("InvestInLoan", mock.Anything, 1, mock.AnythingOfType("*models.LoanInvestment")).Return(nil)

	handler.InvestInLoan(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockLoanService.AssertExpectations(t)
}

func TestLoanHandlerDisburseLoan(t *testing.T) {
	mockLoanService := mocks.NewLoanService(t)
	mockEmailService := mocks2.NewEmailService(t)
	mockStorageService := mocks2.NewStorageService(t)

	handler := NewLoanHandler(mockLoanService, mockEmailService, mockStorageService)

	disbursementReq := &models.LoanDisbursement{
		FieldOfficerEmployeeID:   "emp002",
		AgreementLetterSignedUrl: "https://example.com/signed-agreement.pdf",
	}

	disbursementReqBytes, _ := json.Marshal(disbursementReq)

	req, _ := http.NewRequest("POST", "/api/v1/loans/1/disburse", bytes.NewBuffer(disbursementReqBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Set up chi URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockLoanService.On("DisburseLoan", mock.Anything, 1, mock.AnythingOfType("*models.LoanDisbursement")).Return(nil)

	handler.DisburseLoan(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockLoanService.AssertExpectations(t)
}