package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kitabisa/loan-engine/internal/models"
	"github.com/kitabisa/loan-engine/internal/services"
	"github.com/kitabisa/loan-engine/pkg/external"

	"github.com/go-chi/chi/v5"
)

type LoanHandler struct {
	loanService      services.LoanService
	emailService     external.EmailService
	storageService   external.StorageService
}

func NewLoanHandler(loanService services.LoanService, emailService external.EmailService, storageService external.StorageService) *LoanHandler {
	return &LoanHandler{
		loanService:    loanService,
		emailService:   emailService,
		storageService: storageService,
	}
}

func (h *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var loan struct {
		BorrowerID          int     `json:"borrower_id"`
		PrincipalAmount     float64 `json:"principal_amount"`
		Rate                float64 `json:"rate"`
		ROI                 float64 `json:"roi"`
		AgreementLetterLink string  `json:"agreement_letter_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.Loan{
		BorrowerID:          loan.BorrowerID,
		PrincipalAmount:     loan.PrincipalAmount,
		Rate:                loan.Rate,
		ROI:                 loan.ROI,
		AgreementLetterLink: getNullString(loan.AgreementLetterLink),
	}

	if err := h.loanService.CreateLoan(r.Context(), model); err != nil {
		SendErrorResponse(w, "Failed to create loan", err)
		return
	}

	SendSuccessResponse(w, model, "Loan created successfully")
}

func (h *LoanHandler) GetLoanByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid loan ID", err)
		return
	}

	loan, err := h.loanService.GetLoanByID(r.Context(), id)
	if err != nil {
		SendErrorResponse(w, "Failed to get loan", err)
		return
	}

	SendSuccessResponse(w, loan, "Loan retrieved successfully")
}

func (h *LoanHandler) UpdateLoan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid loan ID", err)
		return
	}

	var loan struct {
		BorrowerID          int     `json:"borrower_id"`
		PrincipalAmount     float64 `json:"principal_amount"`
		Rate                float64 `json:"rate"`
		ROI                 float64 `json:"roi"`
		AgreementLetterLink string  `json:"agreement_letter_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loan); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.Loan{
		BorrowerID:          loan.BorrowerID,
		PrincipalAmount:     loan.PrincipalAmount,
		Rate:                loan.Rate,
		ROI:                 loan.ROI,
		AgreementLetterLink: getNullString(loan.AgreementLetterLink),
	}

	if err := h.loanService.UpdateLoan(r.Context(), id, model); err != nil {
		SendErrorResponse(w, "Failed to update loan", err)
		return
	}

	SendSuccessResponse(w, model, "Loan updated successfully")
}

func (h *LoanHandler) DeleteLoan(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid loan ID", err)
		return
	}

	if err := h.loanService.DeleteLoan(r.Context(), id); err != nil {
		SendErrorResponse(w, "Failed to delete loan", err)
		return
	}

	SendSuccessResponse(w, nil, "Loan deleted successfully")
}

func (h *LoanHandler) ListLoans(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		state = ""
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	var statePtr *string
	if state != "" {
		statePtr = &state
	}

	loans, err := h.loanService.ListLoans(r.Context(), statePtr, offset, limit)
	if err != nil {
		SendErrorResponse(w, "Failed to list loans", err)
		return
	}

	SendSuccessResponse(w, loans, "Loans retrieved successfully")
}

func (h *LoanHandler) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid loan ID", err)
		return
	}

	var approvalData struct {
		FieldValidatorEmployeeID string `json:"field_validator_employee_id"`
		ProofImageUrl            string `json:"proof_image_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&approvalData); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.LoanApproval{
		FieldValidatorEmployeeID: approvalData.FieldValidatorEmployeeID,
		ProofImageUrl:            approvalData.ProofImageUrl,
	}

	if err := h.loanService.ApproveLoan(r.Context(), loanID, model); err != nil {
		SendErrorResponse(w, "Failed to approve loan", err)
		return
	}

	SendSuccessResponse(w, nil, "Loan approved successfully")
}

func (h *LoanHandler) InvestInLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid loan ID", err)
		return
	}

	var investmentData struct {
		InvestorID       int     `json:"investor_id"`
		InvestmentAmount float64 `json:"investment_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&investmentData); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.LoanInvestment{
		InvestorID:       investmentData.InvestorID,
		InvestmentAmount: investmentData.InvestmentAmount,
	}

	if err := h.loanService.InvestInLoan(r.Context(), loanID, model); err != nil {
		SendErrorResponse(w, "Failed to invest in loan", err)
		return
	}

	SendSuccessResponse(w, nil, "Investment completed successfully")
}

func (h *LoanHandler) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	loanID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid loan ID", err)
		return
	}

	var disbursementData struct {
		FieldOfficerEmployeeID      string `json:"field_officer_employee_id"`
		AgreementLetterSignedUrl    string `json:"agreement_letter_signed_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&disbursementData); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.LoanDisbursement{
		FieldOfficerEmployeeID:      disbursementData.FieldOfficerEmployeeID,
		AgreementLetterSignedUrl:    disbursementData.AgreementLetterSignedUrl,
	}

	if err := h.loanService.DisburseLoan(r.Context(), loanID, model); err != nil {
		SendErrorResponse(w, "Failed to disburse loan", err)
		return
	}

	SendSuccessResponse(w, nil, "Loan disbursed successfully")
}

func (h *LoanHandler) GetLoansByState(w http.ResponseWriter, r *http.Request) {
	state := chi.URLParam(r, "state")

	loans, err := h.loanService.GetLoansByState(r.Context(), state)
	if err != nil {
		SendErrorResponse(w, "Failed to get loans by state", err)
		return
	}

	SendSuccessResponse(w, loans, "Loans retrieved successfully")
}

// Helper function to convert string to sql.NullString
func getNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}