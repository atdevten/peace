package handlers

import (
	"github.com/atdevten/peace/internal/application/usecases"
	"github.com/atdevten/peace/internal/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecases.UserUseCase
}

func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromGinContext(c)
	if !ok {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}
	ctx := c.Request.Context()
	user, err := h.userUseCase.GetByID(ctx, userID.String())
	if err != nil {
		Error(c, CodeNotFound, "User not found")
		return
	}
	data := gin.H{
		"id":        user.ID().String(),
		"email":     user.Email().String(),
		"username":  user.Username().String(),
		"full_name": user.GetFullName(),
	}
	Success(c, "Me retrieved successfully", data)
}

type UpdateProfileRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromGinContext(c)
	if !ok {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, CodeBadRequest, "Invalid request body")
		return
	}
	ctx := c.Request.Context()
	user, err := h.userUseCase.UpdateProfile(ctx, userID.String(), req.FirstName, req.LastName)
	if err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}
	data := gin.H{
		"id":        user.ID().String(),
		"email":     user.Email().String(),
		"username":  user.Username().String(),
		"full_name": user.GetFullName(),
	}
	Success(c, "Profile updated successfully", data)
}

type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromGinContext(c)
	if !ok {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.NewPassword == "" {
		Error(c, CodeBadRequest, "new_password is required")
		return
	}
	ctx := c.Request.Context()
	if err := h.userUseCase.UpdatePassword(ctx, userID.String(), req.NewPassword); err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}
	Success(c, "Password updated successfully", nil)
}

func (h *UserHandler) Deactivate(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromGinContext(c)
	if !ok {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}
	ctx := c.Request.Context()
	if err := h.userUseCase.Deactivate(ctx, userID.String()); err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}
	Success(c, "Account deactivated successfully", nil)
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, ok := middleware.GetUserIDFromGinContext(c)
	if !ok {
		Error(c, CodeUnauthorized, "User not authenticated")
		return
	}
	ctx := c.Request.Context()
	if err := h.userUseCase.Delete(ctx, userID.String()); err != nil {
		Error(c, CodeBadRequest, err.Error())
		return
	}
	Success(c, "Account deleted successfully", nil)
}
