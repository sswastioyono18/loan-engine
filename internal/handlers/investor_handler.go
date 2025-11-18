package handlers

import (
	"encoding/json"
	"github.com/sswastioyono18/loan-engine/internal/models"
	"github.com/sswastioyono18/loan-engine/internal/services"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type InvestorHandler struct {
	investorService services.InvestorService
}

func NewInvestorHandler(investorService services.InvestorService) *InvestorHandler {
	return &InvestorHandler{
		investorService: investorService,
	}
}

func (h *InvestorHandler) CreateInvestor(w http.ResponseWriter, r *http.Request) {
	var investor struct {
		InvestorID string `json:"investor_id"`
		FullName   string `json:"full_name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&investor); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.Investor{
		InvestorID: investor.InvestorID,
		FullName:   investor.FullName,
		Email:      investor.Email,
		Phone:      investor.Phone,
	}

	if err := h.investorService.CreateInvestor(r.Context(), model); err != nil {
		SendErrorResponse(w, "Failed to create investor", err)
		return
	}

	SendSuccessResponse(w, model, "Investor created successfully")
}

func (h *InvestorHandler) GetInvestorByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid investor ID", err)
		return
	}

	investor, err := h.investorService.GetInvestorByID(r.Context(), id)
	if err != nil {
		SendErrorResponse(w, "Failed to get investor", err)
		return
	}

	SendSuccessResponse(w, investor, "Investor retrieved successfully")
}

func (h *InvestorHandler) UpdateInvestor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid investor ID", err)
		return
	}

	var investor struct {
		InvestorID string `json:"investor_id"`
		FullName   string `json:"full_name"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&investor); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.Investor{
		ID:         id,
		InvestorID: investor.InvestorID,
		FullName:   investor.FullName,
		Email:      investor.Email,
		Phone:      investor.Phone,
	}

	if err := h.investorService.UpdateInvestor(r.Context(), id, model); err != nil {
		SendErrorResponse(w, "Failed to update investor", err)
		return
	}

	SendSuccessResponse(w, model, "Investor updated successfully")
}

func (h *InvestorHandler) DeleteInvestor(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid investor ID", err)
		return
	}

	if err := h.investorService.DeleteInvestor(r.Context(), id); err != nil {
		SendErrorResponse(w, "Failed to delete investor", err)
		return
	}

	SendSuccessResponse(w, nil, "Investor deleted successfully")
}

func (h *InvestorHandler) ListInvestors(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	investors, err := h.investorService.ListInvestors(r.Context(), offset, limit)
	if err != nil {
		SendErrorResponse(w, "Failed to list investors", err)
		return
	}

	SendSuccessResponse(w, investors, "Investors retrieved successfully")
}
