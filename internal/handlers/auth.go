package handlers

import (
	"net/http"

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
	utils.GinSuccessResponse(c, "Users retrieved", []interface{}{})
}

func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	utils.GinSuccessResponse(c, "User role updated", nil)
}
