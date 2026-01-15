package handlers

import (
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	// Parse UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	// Get cart from service
	cart, err := h.cartService.GetCart(c.Request.Context(), userUUID)
	if err != nil {
		utils.GinNotFoundResponse(c, "Cart")
		return
	}

	utils.GinSuccessResponse(c, "Cart retrieved successfully", cart)
}

func (h *CartHandler) ValidateCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	// Parse UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	// Get cart to validate
	cart, err := h.cartService.GetCart(c.Request.Context(), userUUID)
	if err != nil {
		utils.GinNotFoundResponse(c, "Cart")
		return
	}

	// Get cart ID from the cart
	cartID := cart.ID

	// Validate cart via service
	isValid, errors, err := h.cartService.ValidateCart(c.Request.Context(), cartID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to validate cart", err)
		return
	}

	if !isValid {
		c.JSON(409, map[string]interface{}{
			"success": false,
			"message": "Cart validation failed",
			"data": map[string]interface{}{
				"valid":  false,
				"errors": errors,
			},
		})
		return
	}

	utils.GinSuccessResponse(c, "Cart is valid", map[string]interface{}{
		"valid": true,
		"cart":  cart,
	})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	// Parse UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	// Add to cart via service
	cart, err := h.cartService.AddToCart(c.Request.Context(), userUUID, req)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to add item to cart", err)
		return
	}

	utils.GinCreatedResponse(c, "Item added to cart successfully", cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	itemID := c.Param("itemId")
	if itemID == "" {
		utils.GinBadRequestResponse(c, "Item ID is required", nil)
		return
	}

	var req models.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	// Parse UUIDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid item ID", err)
		return
	}

	// Update cart item via service
	cart, err := h.cartService.UpdateCartItem(c.Request.Context(), userUUID, itemUUID, req)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to update cart item", err)
		return
	}

	utils.GinSuccessResponse(c, "Cart item updated successfully", cart)
}

func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	itemID := c.Param("itemId")
	if itemID == "" {
		utils.GinBadRequestResponse(c, "Item ID is required", nil)
		return
	}

	// Parse UUIDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid item ID", err)
		return
	}

	// Remove from cart via service
	cart, err := h.cartService.RemoveFromCart(c.Request.Context(), userUUID, itemUUID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to remove item from cart", err)
		return
	}

	utils.GinSuccessResponse(c, "Item removed from cart successfully", cart)
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	// Parse UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	// Clear cart via service
	err = h.cartService.ClearCart(c.Request.Context(), userUUID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to clear cart", err)
		return
	}

	utils.GinSuccessResponse(c, "Cart cleared successfully", nil)
}
