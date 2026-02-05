package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	AuthService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	// Register user
	user, err := h.AuthService.Register(c.Request.Context(), req)
	if err != nil {
		utils.GinBadRequestResponse(c, "Registration failed", err)
		return
	}

	// Generate token
	token, err := h.AuthService.GenerateToken(user)
	if err != nil {
		utils.GinInternalErrorResponse(c, "Failed to generate token", err)
		return
	}

	response := models.LoginResponse{
		User:        user,
		AccessToken: token,
	}

	utils.GinCreatedResponse(c, "User registered successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	// Login user
	response, err := h.AuthService.Login(c.Request.Context(), req)
	if err != nil {
		utils.GinErrorResponse(c, http.StatusUnauthorized, "Login failed", err)
		return
	}

	utils.GinSuccessResponse(c, "Login successful", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	user, err := h.AuthService.ValidateToken(req.RefreshToken)
	if err != nil {
		utils.GinUnauthorizedResponse(c, "Invalid refresh token")
		return
	}

	// Load full user to ensure current role/email
	fullUser, err := h.AuthService.GetProfile(c.Request.Context(), user.ID)
	if err != nil {
		utils.GinUnauthorizedResponse(c, "Invalid refresh token")
		return
	}

	token, err := h.AuthService.GenerateToken(fullUser)
	if err != nil {
		utils.GinInternalErrorResponse(c, "Failed to generate token", err)
		return
	}

	utils.GinSuccessResponse(c, "Token refreshed", gin.H{"access_token": token})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   err.Error(),
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
		return
	}

	// Get user profile
	user, err := h.AuthService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
			"error":   err.Error(),
		})
		return
	}

	utils.GinSuccessResponse(c, "Profile retrieved successfully", user)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// Get user ID from context
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
			"error":   err.Error(),
		})
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email" validate:"email"`
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); errors != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    nil,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=6"`
	}

	// Parse request body
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	utils.GinSuccessResponse(c, "Password changed successfully", nil)
}

// Stub methods for admin endpoints
func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	rangeDays := 0
	if rd := c.Query("range_days"); rd != "" {
		if parsed, err := strconv.Atoi(rd); err == nil && parsed > 0 {
			rangeDays = parsed
		}
	}

	users, total, err := h.AuthService.ListUsers(c.Request.Context(), page, limit, rangeDays)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to retrieve users", err)
		return
	}

	response := map[string]interface{}{
		"users": users,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	}

	utils.GinSuccessResponse(c, "Users retrieved", response)
}

func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	var req models.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	updatedUser, token, err := h.AuthService.UpdateUserRole(c.Request.Context(), userUUID, req.Role)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to update user role", err)
		return
	}

	response := gin.H{
		"message": "User role updated",
		"user":    updatedUser,
		"token":   token,
	}

	utils.GinSuccessResponse(c, "User role updated successfully", response)
}
