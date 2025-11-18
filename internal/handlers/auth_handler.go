package handlers

import (
	"encoding/json"
	"github.com/sswastioyono18/loan-engine/internal/models"
	"github.com/sswastioyono18/loan-engine/internal/services"
	"net/http"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		UserID   string `json:"user_id"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"user_type"`
		FullName string `json:"full_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	model := &models.User{
		UserID:   user.UserID,
		Email:    user.Email,
		UserType: user.UserType,
		FullName: user.FullName,
		IsActive: true, // Default to active
	}

	if err := h.authService.RegisterUser(r.Context(), model, user.Password); err != nil {
		SendErrorResponse(w, "Failed to register user", err)
		return
	}

	SendSuccessResponse(w, nil, "User registered successfully")
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	token, err := h.authService.LoginUser(r.Context(), credentials.Email, credentials.Password)
	if err != nil {
		SendErrorResponse(w, "Login failed", err)
		return
	}

	SendSuccessResponse(w, map[string]string{"token": token}, "Login successful")
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, "Invalid request body", err)
		return
	}

	newToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		SendErrorResponse(w, "Token refresh failed", err)
		return
	}

	SendSuccessResponse(w, map[string]string{"token": newToken}, "Token refreshed successfully")
}
