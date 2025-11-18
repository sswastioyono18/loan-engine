package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sswastioyono18/loan-engine/internal/models"
	"github.com/sswastioyono18/loan-engine/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, user *models.User, password string) error
	LoginUser(ctx context.Context, email, password string) (string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type authServiceImpl struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string) AuthService {
	return &authServiceImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authServiceImpl) RegisterUser(ctx context.Context, user *models.User, password string) error {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Hash the password
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = hashedPassword
	user.IsActive = true

	return s.userRepo.Create(ctx, user)
}

func (s *authServiceImpl) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return "", fmt.Errorf("user account is deactivated")
	}

	if !s.CheckPasswordHash(password, user.PasswordHash) {
		return "", fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:   user.ID,
		Email:    user.Email,
		UserType: user.UserType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "loan-engine",
		},
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return tokenString, nil
}

func (s *authServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// In a real implementation, you would validate the refresh token
	// For now, we'll just generate a new access token
	// This is a simplified implementation - in production, you'd want to store and validate refresh tokens

	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Generate new access token
		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			UserID:   claims.UserID,
			Email:    claims.Email,
			UserType: claims.UserType,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				Issuer:    "loan-engine",
			},
		})

		newTokenString, err := newToken.SignedString([]byte(s.jwtSecret))
		if err != nil {
			return "", fmt.Errorf("failed to generate new token: %w", err)
		}

		return newTokenString, nil
	}

	return "", fmt.Errorf("invalid refresh token")
}

func (s *authServiceImpl) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		user, err := s.userRepo.GetByID(ctx, claims.UserID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *authServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *authServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
