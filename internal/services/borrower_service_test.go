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

func TestCreateBorrower(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	borrower := &models.Borrower{
		BorrowerIDNumber: "B001",
		FullName:         "John Doe",
		Email:            "john@example.com",
		Phone:            "1234567890",
		Address:          "123 Main St",
	}

	// Test successful creation
	mockRepo.On("Create", context.Background(), borrower).Return(nil)

	err := service.CreateBorrower(context.Background(), borrower)

	assert.NoError(t, err)
}

func TestCreateBorrowerError(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	borrower := &models.Borrower{
		BorrowerIDNumber: "B001",
		FullName:         "John Doe",
		Email:            "john@example.com",
		Phone:            "1234567890",
		Address:          "123 Main St",
	}

	// Test creation error
	mockRepo.On("Create", context.Background(), borrower).Return(errors.New("database error"))

	err := service.CreateBorrower(context.Background(), borrower)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestGetBorrowerByID(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	borrower := &models.Borrower{
		ID:               1,
		BorrowerIDNumber: "B001",
		FullName:         "John Doe",
		Email:            "john@example.com",
		Phone:            "1234567890",
		Address:          "123 Main St",
	}

	// Test successful retrieval
	mockRepo.On("GetByID", context.Background(), 1).Return(borrower, nil)

	result, err := service.GetBorrowerByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, borrower, result)
}

func TestGetBorrowerByIDNotFound(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	// Test not found
	mockRepo.On("GetByID", context.Background(), 1).Return(nil, errors.New("borrower not found"))

	_, err := service.GetBorrowerByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "borrower not found")
}

func TestGetBorrowerByBorrowerIDNumber(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	borrower := &models.Borrower{
		ID:               1,
		BorrowerIDNumber: "B001",
		FullName:         "John Doe",
		Email:            "john@example.com",
		Phone:            "1234567890",
		Address:          "123 Main St",
	}

	// Test successful retrieval by ID number
	mockRepo.On("GetByBorrowerIDNumber", context.Background(), "B001").Return(borrower, nil)

	result, err := service.GetBorrowerByBorrowerIDNumber(context.Background(), "B001")

	assert.NoError(t, err)
	assert.Equal(t, borrower, result)
}

func TestGetBorrowerByBorrowerIDNumberNotFound(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	// Test not found by ID number
	mockRepo.On("GetByBorrowerIDNumber", context.Background(), "B001").Return(nil, errors.New("borrower not found"))

	_, err := service.GetBorrowerByBorrowerIDNumber(context.Background(), "B001")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "borrower not found")
}

func TestUpdateBorrower(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	existingBorrower := &models.Borrower{
		ID:               1,
		BorrowerIDNumber: "B001",
		FullName:         "John Doe",
		Email:            "john@example.com",
		Phone:            "1234567890",
		Address:          "123 Main St",
	}

	updatedBorrower := &models.Borrower{
		BorrowerIDNumber: "B002",
		FullName:         "Jane Doe",
		Email:            "jane@example.com",
		Phone:            "0987654321",
		Address:          "456 Oak Ave",
	}

	// Test successful update
	mockRepo.On("GetByID", context.Background(), 1).Return(existingBorrower, nil)
	mockRepo.On("Update", context.Background(), mock.AnythingOfType("*models.Borrower")).Return(nil)

	err := service.UpdateBorrower(context.Background(), 1, updatedBorrower)

	assert.NoError(t, err)
	// Verify that the ID and CreatedAt were preserved from the existing borrower
	assert.Equal(t, 1, updatedBorrower.ID)
	assert.Equal(t, existingBorrower.CreatedAt, updatedBorrower.CreatedAt)
}

func TestUpdateBorrowerNotFound(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	updatedBorrower := &models.Borrower{
		BorrowerIDNumber: "B002",
		FullName:         "Jane Doe",
		Email:            "jane@example.com",
		Phone:            "0987654321",
		Address:          "456 Oak Ave",
	}

	// Test update when borrower doesn't exist
	mockRepo.On("GetByID", context.Background(), 1).Return(nil, errors.New("borrower not found"))

	err := service.UpdateBorrower(context.Background(), 1, updatedBorrower)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "borrower not found")
}

func TestUpdateBorrowerUpdateError(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	existingBorrower := &models.Borrower{
		ID:               1,
		BorrowerIDNumber: "B001",
		FullName:         "John Doe",
		Email:            "john@example.com",
		Phone:            "1234567890",
		Address:          "123 Main St",
	}

	updatedBorrower := &models.Borrower{
		ID:               1,
		BorrowerIDNumber: "B002",
		FullName:         "Jane Doe",
		Email:            "jane@example.com",
		Phone:            "0987654321",
		Address:          "456 Oak Ave",
	}

	// Test update error
	mockRepo.On("GetByID", context.Background(), 1).Return(existingBorrower, nil)
	mockRepo.On("Update", context.Background(), updatedBorrower).Return(errors.New("update failed"))

	err := service.UpdateBorrower(context.Background(), 1, updatedBorrower)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update failed")
}

func TestDeleteBorrower(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	// Test successful deletion
	mockRepo.On("Delete", context.Background(), 1).Return(nil)

	err := service.DeleteBorrower(context.Background(), 1)

	assert.NoError(t, err)
}

func TestDeleteBorrowerError(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	// Test deletion error
	mockRepo.On("Delete", context.Background(), 1).Return(errors.New("delete failed"))

	err := service.DeleteBorrower(context.Background(), 1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete failed")
}

func TestListBorrowers(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	borrowers := []*models.Borrower{
		{
			ID:               1,
			BorrowerIDNumber: "B001",
			FullName:         "John Doe",
			Email:            "john@example.com",
			Phone:            "1234567890",
			Address:          "123 Main St",
		},
		{
			ID:               2,
			BorrowerIDNumber: "B002",
			FullName:         "Jane Doe",
			Email:            "jane@example.com",
			Phone:            "0987654321",
			Address:          "456 Oak Ave",
		},
	}

	// Test successful listing
	mockRepo.On("List", context.Background(), 0, 10).Return(borrowers, nil)

	result, err := service.ListBorrowers(context.Background(), 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, borrowers, result)
}

func TestListBorrowersError(t *testing.T) {
	mockRepo := mocks.NewBorrowerRepository(t)
	service := NewBorrowerService(mockRepo)

	// Test listing error
	mockRepo.On("List", context.Background(), 0, 10).Return(nil, errors.New("list failed"))

	_, err := service.ListBorrowers(context.Background(), 0, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "list failed")
}
