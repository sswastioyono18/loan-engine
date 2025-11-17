package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/kitabisa/loan-engine/internal/models"
	"github.com/kitabisa/loan-engine/internal/services"

	"github.com/go-chi/chi/v5"
)

type BorrowerHandler struct {
	borrowerService services.BorrowerService
}

func NewBorrowerHandler(borrowerService services.BorrowerService) *BorrowerHandler {
	return &BorrowerHandler{
		borrowerService: borrowerService,
	}
}

func (h *BorrowerHandler) CreateBorrower(w http.ResponseWriter, r *http.Request) {
	var borrower struct {
		BorrowerIDNumber string `json:"borrower_id_number"`
		FullName         string `json:"full_name"`
		Email            string `json:"email"`
		Phone            string `json:"phone"`
		Address          string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&borrower); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.Borrower{
		BorrowerIDNumber: borrower.BorrowerIDNumber,
		FullName:         borrower.FullName,
		Email:            borrower.Email,
		Phone:            borrower.Phone,
		Address:          borrower.Address,
	}

	if err := h.borrowerService.CreateBorrower(r.Context(), model); err != nil {
		SendErrorResponse(w, "Failed to create borrower", err)
		return
	}

	SendSuccessResponse(w, model, "Borrower created successfully")
}

func (h *BorrowerHandler) GetBorrowerByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid borrower ID", err)
		return
	}

	borrower, err := h.borrowerService.GetBorrowerByID(r.Context(), id)
	if err != nil {
		SendErrorResponse(w, "Failed to get borrower", err)
		return
	}

	SendSuccessResponse(w, borrower, "Borrower retrieved successfully")
}

func (h *BorrowerHandler) UpdateBorrower(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid borrower ID", err)
		return
	}

	var borrower struct {
		BorrowerIDNumber string `json:"borrower_id_number"`
		FullName         string `json:"full_name"`
		Email            string `json:"email"`
		Phone            string `json:"phone"`
		Address          string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&borrower); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.Borrower{
		ID:               id,
		BorrowerIDNumber: borrower.BorrowerIDNumber,
		FullName:         borrower.FullName,
		Email:            borrower.Email,
		Phone:            borrower.Phone,
		Address:          borrower.Address,
	}

	if err := h.borrowerService.UpdateBorrower(r.Context(), id, model); err != nil {
		SendErrorResponse(w, "Failed to update borrower", err)
		return
	}

	SendSuccessResponse(w, model, "Borrower updated successfully")
}

func (h *BorrowerHandler) DeleteBorrower(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorResponse(w, "Invalid borrower ID", err)
		return
	}

	if err := h.borrowerService.DeleteBorrower(r.Context(), id); err != nil {
		SendErrorResponse(w, "Failed to delete borrower", err)
		return
	}

	SendSuccessResponse(w, nil, "Borrower deleted successfully")
}

func (h *BorrowerHandler) ListBorrowers(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	borrowers, err := h.borrowerService.ListBorrowers(r.Context(), offset, limit)
	if err != nil {
		SendErrorResponse(w, "Failed to list borrowers", err)
		return
	}

	SendSuccessResponse(w, borrowers, "Borrowers retrieved successfully")
}