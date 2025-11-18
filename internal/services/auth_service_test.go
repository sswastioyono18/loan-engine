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

func TestRegisterUser(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		Email:    "test@example.com",
		UserType: "investor",
		FullName: "Test User",
	}

	// Test successful registration
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(nil, errors.New("user not found"))
	mockUserRepo.On("Create", context.Background(), mock.AnythingOfType("*models.User")).Return(nil)

	err := service.RegisterUser(context.Background(), user, "password123")

	assert.NoError(t, err)
	assert.True(t, len(user.PasswordHash) > 0) // Password should be hashed
	assert.True(t, user.IsActive)
}

func TestRegisterUserDuplicateEmail(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		Email:    "test@example.com",
		UserType: "investor",
		FullName: "Test User",
	}

	existingUser := &models.User{
		ID:       1,
		Email:    "test@example.com",
		UserType: "investor",
		FullName: "Existing User",
	}

	// Test duplicate email
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(existingUser, nil)

	err := service.RegisterUser(context.Background(), user, "password123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestRegisterUserPasswordHashError(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		Email:    "test@example.com",
		UserType: "investor",
		FullName: "Test User",
	}

	// Test successful registration with empty password (should still work)
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(nil, errors.New("user not found"))
	mockUserRepo.On("Create", context.Background(), mock.AnythingOfType("*models.User")).Return(nil)

	err := service.RegisterUser(context.Background(), user, "") // Empty password

	assert.NoError(t, err)
	assert.True(t, len(user.PasswordHash) > 0) // Password should be hashed even if empty
}

func TestLoginUser(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		ID:           1,
		Email:        "test@example.com",
		UserType:     "investor",
		FullName:     "Test User",
		PasswordHash: "$2a$14$qxXQWcJG23rX0daSNJl6FO8I4V9Hj55ibaqUqzZHaa7x0UXv2djLa", // bcrypt hash for "password123"
		IsActive:     true,
	}

	// Test successful login
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(user, nil)

	token, err := service.LoginUser(context.Background(), user.Email, "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	// Test invalid email
	mockUserRepo.On("GetByEmail", context.Background(), "nonexistent@example.com").Return(nil, errors.New("user not found"))

	_, err := service.LoginUser(context.Background(), "nonexistent@example.com", "password123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestLoginUserInactiveAccount(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		ID:           1,
		Email:        "test@example.com",
		UserType:     "investor",
		FullName:     "Test User",
		PasswordHash: "$2a$14$qxXQWcJG23rX0daSNJl6FO8I4V9Hj55ibaqUqzZHaa7x0UXv2djLa", // bcrypt hash for "password123"
		IsActive:     false,
	}

	// Test inactive account
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(user, nil)

	_, err := service.LoginUser(context.Background(), user.Email, "password123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user account is deactivated")
}

func TestLoginUserInvalidPassword(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		ID:           1,
		Email:        "test@example.com",
		UserType:     "investor",
		FullName:     "Test User",
		PasswordHash: "$2a$14$qxXQWcJG23rX0daSNJl6FO8I4V9Hj55ibaqUqzZHaa7x0UXv2djLa", // bcrypt hash for "password123"
		IsActive:     true,
	}

	// Test invalid password
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(user, nil)

	_, err := service.LoginUser(context.Background(), user.Email, "wrongpassword")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestValidateToken(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		UserType: "investor",
		FullName: "Test User",
		IsActive: true,
		PasswordHash: "$2a$14$qxXQWcJG23rX0daSNJl6FO8I4V9Hj55ibaqUqzZHaa7x0UXv2djLa", // bcrypt hash for "password123"
	}

	// Create a valid token first
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(user, nil)
	token, err := service.LoginUser(context.Background(), user.Email, "password123")
	assert.NoError(t, err)

	// Now test token validation
	mockUserRepo.On("GetByID", context.Background(), user.ID).Return(user, nil)

	validatedUser, err := service.ValidateToken(context.Background(), token)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, validatedUser.ID)
	assert.Equal(t, user.Email, validatedUser.Email)
}

func TestValidateTokenInvalid(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	// Test invalid token
	_, err := service.ValidateToken(context.Background(), "invalid-token")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestValidateTokenUserNotFound(t *testing.T) {
	mockUserRepo := mocks.NewUserRepository(t)
	service := NewAuthService(mockUserRepo, "test-secret")

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		UserType: "investor",
		FullName: "Test User",
		IsActive: true,
		PasswordHash: "$2a$14$qxXQWcJG23rX0daSNJl6FO8I4V9Hj55ibaqUqzZHaa7x0UXv2djLa", // bcrypt hash for "password123"
	}

	// Create a valid token first
	mockUserRepo.On("GetByEmail", context.Background(), user.Email).Return(user, nil)
	token, err := service.LoginUser(context.Background(), user.Email, "password123")
	assert.NoError(t, err)

	// Now test token validation with user not found
	mockUserRepo.On("GetByID", context.Background(), user.ID).Return(nil, errors.New("user not found"))

	_, err = service.ValidateToken(context.Background(), token)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestHashPassword(t *testing.T) {
	service := NewAuthService(nil, "test-secret")

	hash, err := service.HashPassword("password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.True(t, len(hash) > 10) // bcrypt hash should be long
}

func TestCheckPasswordHash(t *testing.T) {
	service := NewAuthService(nil, "test-secret")

	password := "password123"
	hash, err := service.HashPassword(password)
	assert.NoError(t, err)

	// Test correct password
	result := service.CheckPasswordHash(password, hash)
	assert.True(t, result)

	// Test incorrect password
	result = service.CheckPasswordHash("wrongpassword", hash)
	assert.False(t, result)
}
