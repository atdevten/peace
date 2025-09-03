package handlers

import (
	"github.com/atdevten/peace/internal/application/commands"
	"github.com/atdevten/peace/internal/application/usecases"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase usecases.AuthUseCase
}

type RegisterRequest struct {
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
}

type LoginResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

func NewAuthHandler(authUseCase usecases.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	// Create command
	command, err := commands.NewRegisterCommand(
		req.Email,
		req.Username,
		req.Password,
		req.FirstName,
		req.LastName,
	)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	user, err := h.authUseCase.Register(ctx, command)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Build response
	userResponse := UserResponse{
		ID:       user.ID().String(),
		Email:    user.Email().String(),
		Username: user.Username().String(),
		FullName: user.GetFullName(),
	}

	c.JSON(201, APIResponse{
		Code:    CodeSuccess,
		Message: "User registered successfully",
		Data:    userResponse,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	// Create command
	command, err := commands.NewLoginCommand(req.Email, req.Password)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	user, access, refresh, err := h.authUseCase.Login(ctx, command)
	if err != nil {
		Error(c, CodeUnauthorized, err.Error())
		return
	}

	// Build response
	userResponse := UserResponse{
		ID:       user.ID().String(),
		Email:    user.Email().String(),
		Username: user.Username().String(),
		FullName: user.GetFullName(),
	}

	loginResponse := LoginResponse{
		User:         userResponse,
		AccessToken:  access,
		RefreshToken: refresh,
	}

	Success(c, "Login successful", loginResponse)
}

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	ctx := c.Request.Context()
	newAccess, newRefresh, err := h.authUseCase.Refresh(ctx, req.AccessToken, req.RefreshToken)
	if err != nil {
		Error(c, CodeUnauthorized, err.Error())
		return
	}

	c.JSON(200, APIResponse{
		Code:    CodeSuccess,
		Message: "Token refreshed successfully",
		Data: gin.H{
			"access_token":  newAccess,
			"refresh_token": newRefresh,
		},
	})
}

// GoogleLoginRequest represents the request for Google OAuth login
type GoogleLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// GoogleLogin handles Google OAuth login
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}

	// Execute use case
	ctx := c.Request.Context()
	user, access, refresh, err := h.authUseCase.LoginWithGoogle(ctx, req.Code)
	if err != nil {
		Error(c, CodeUnauthorized, err.Error())
		return
	}

	// Build response
	userResponse := UserResponse{
		ID:       user.ID().String(),
		Email:    user.Email().String(),
		Username: user.Username().String(),
		FullName: user.GetFullName(),
	}

	loginResponse := LoginResponse{
		User:         userResponse,
		AccessToken:  access,
		RefreshToken: refresh,
	}

	Success(c, "Google login successful", loginResponse)
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (h *AuthHandler) GetGoogleAuthURL(c *gin.Context) {
	// This would typically get the URL from a service
	// For now, we'll return a placeholder
	Success(c, "Google OAuth URL", gin.H{
		"auth_url": "https://accounts.google.com/o/oauth2/v2/auth?client_id=...",
		"note":     "This endpoint should be updated to use Google service",
	})
}
