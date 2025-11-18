package services

import (
	"context"
	"errors"
	"testing"

	"github.com/sswastioyono18/loan-engine/internal/models"
	"github.com/sswastioyono18/loan-engine/internal/repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateInvestor(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	investor := &models.Investor{
		InvestorID: "INV001",
		FullName:   "John Investor",
		Email:      "john@example.com",
		Phone:      "1234567890",
	}

	// Test successful creation
	mockRepo.On("Create", context.Background(), investor).Return(nil)

	err := service.CreateInvestor(context.Background(), investor)

	assert.NoError(t, err)
}

func TestCreateInvestorError(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	investor := &models.Investor{
		InvestorID: "INV001",
		FullName:   "John Investor",
		Email:      "john@example.com",
		Phone:      "1234567890",
	}

	// Test creation error
	mockRepo.On("Create", context.Background(), investor).Return(errors.New("database error"))

	err := service.CreateInvestor(context.Background(), investor)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestGetInvestorByID(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	investor := &models.Investor{
		ID:         1,
		InvestorID: "INV001",
		FullName:   "John Investor",
		Email:      "john@example.com",
		Phone:      "1234567890",
	}

	// Test successful retrieval
	mockRepo.On("GetByID", context.Background(), 1).Return(investor, nil)

	result, err := service.GetInvestorByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, investor, result)
}

func TestGetInvestorByIDNotFound(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	// Test not found
	mockRepo.On("GetByID", context.Background(), 1).Return(nil, errors.New("investor not found"))

	_, err := service.GetInvestorByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "investor not found")
}

func TestGetInvestorByInvestorID(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	investor := &models.Investor{
		ID:         1,
		InvestorID: "INV001",
		FullName:   "John Investor",
		Email:      "john@example.com",
		Phone:      "1234567890",
	}

	// Test successful retrieval by investor ID
	mockRepo.On("GetByInvestorID", context.Background(), "INV001").Return(investor, nil)

	result, err := service.GetInvestorByInvestorID(context.Background(), "INV001")

	assert.NoError(t, err)
	assert.Equal(t, investor, result)
}

func TestGetInvestorByInvestorIDNotFound(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	// Test not found by investor ID
	mockRepo.On("GetByInvestorID", context.Background(), "INV001").Return(nil, errors.New("investor not found"))

	_, err := service.GetInvestorByInvestorID(context.Background(), "INV001")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "investor not found")
}

func TestGetInvestorByEmail(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	investor := &models.Investor{
		ID:         1,
		InvestorID: "INV001",
		FullName:   "John Investor",
		Email:      "john@example.com",
		Phone:      "1234567890",
	}

	// Test successful retrieval by email
	mockRepo.On("GetByEmail", context.Background(), "john@example.com").Return(investor, nil)

	result, err := service.GetInvestorByEmail(context.Background(), "john@example.com")

	assert.NoError(t, err)
	assert.Equal(t, investor, result)
}

func TestGetInvestorByEmailNotFound(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	// Test not found by email
	mockRepo.On("GetByEmail", context.Background(), "nonexistent@example.com").Return(nil, errors.New("investor not found"))

	_, err := service.GetInvestorByEmail(context.Background(), "nonexistent@example.com")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "investor not found")
}

func TestUpdateInvestor(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	existingInvestor := &models.Investor{
		ID:         1,
		InvestorID: "INV001",
		FullName:   "John Investor",
		Email:      "john@example.com",
		Phone:      "1234567890",
	}

	updatedInvestor := &models.Investor{
		InvestorID: "INV002",
		FullName:   "Jane Investor",
		Email:      "jane@example.com",
		Phone:      "0987654321",
	}

	// Test successful update
	mockRepo.On("GetByID", context.Background(), 1).Return(existingInvestor, nil)
	mockRepo.On("Update", context.Background(), mock.AnythingOfType("*models.Investor")).Return(nil)

	err := service.UpdateInvestor(context.Background(), 1, updatedInvestor)

	assert.NoError(t, err)
	// Verify that the ID and CreatedAt were preserved from the existing investor
	assert.Equal(t, 1, updatedInvestor.ID)
	assert.Equal(t, existingInvestor.CreatedAt, updatedInvestor.CreatedAt)
}

func TestUpdateInvestorNotFound(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	updatedInvestor := &models.Investor{
		InvestorID: "INV002",
		FullName:   "Jane Investor",
		Email:      "jane@example.com",
		Phone:      "0987654321",
	}

	// Test update when investor doesn't exist
	mockRepo.On("GetByID", context.Background(), 1).Return(nil, errors.New("investor not found"))

	err := service.UpdateInvestor(context.Background(), 1, updatedInvestor)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "investor not found")
}

func TestUpdateInvestorUpdateError(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	existingInvestor := &models.Investor{
		ID:         1,
		InvestorID: "INV001",
		FullName:   "John Investor",
		Email:      "john@example.com",
		Phone:      "1234567890",
	}

	updatedInvestor := &models.Investor{
		ID:         1,
		InvestorID: "INV002",
		FullName:   "Jane Investor",
		Email:      "jane@example.com",
		Phone:      "0987654321",
	}

	// Test update error
	mockRepo.On("GetByID", context.Background(), 1).Return(existingInvestor, nil)
	mockRepo.On("Update", context.Background(), updatedInvestor).Return(errors.New("update failed"))

	err := service.UpdateInvestor(context.Background(), 1, updatedInvestor)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
}

func TestDeleteInvestor(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	// Test successful deletion
	mockRepo.On("Delete", context.Background(), 1).Return(nil)

	err := service.DeleteInvestor(context.Background(), 1)

	assert.NoError(t, err)
}

func TestDeleteInvestorError(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	// Test deletion error
	mockRepo.On("Delete", context.Background(), 1).Return(errors.New("delete failed"))

	err := service.DeleteInvestor(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete failed")
}

func TestListInvestors(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	investors := []*models.Investor{
		{
			ID:         1,
			InvestorID: "INV001",
			FullName:   "John Investor",
			Email:      "john@example.com",
			Phone:      "1234567890",
		},
		{
			ID:         2,
			InvestorID: "INV002",
			FullName:   "Jane Investor",
			Email:      "jane@example.com",
			Phone:      "0987654321",
		},
	}

	// Test successful listing
	mockRepo.On("List", context.Background(), 0, 10).Return(investors, nil)

	result, err := service.ListInvestors(context.Background(), 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, investors, result)
}

func TestListInvestorsError(t *testing.T) {
	mockRepo := mocks.NewInvestorRepository(t)
	service := NewInvestorService(mockRepo)

	// Test listing error
	mockRepo.On("List", context.Background(), 0, 10).Return(nil, errors.New("list failed"))

	_, err := service.ListInvestors(context.Background(), 0, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "list failed")
}
