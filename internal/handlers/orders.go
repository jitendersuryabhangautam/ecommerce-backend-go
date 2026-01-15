package handlers

import (
	"strconv"

	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), userUUID, req)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to create order", err)
		return
	}

	utils.GinCreatedResponse(c, "Order created successfully", order)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	// Get pagination params
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

	orders, total, err := h.orderService.GetUserOrders(c.Request.Context(), userUUID, page, limit)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to retrieve orders", err)
		return
	}

	response := map[string]interface{}{
		"orders": orders,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	}

	utils.GinSuccessResponse(c, "Orders retrieved successfully", response)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	orderID := c.Param("id")
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid order ID", err)
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), orderUUID, userUUID)
	if err != nil {
		utils.GinNotFoundResponse(c, "Order not found")
		return
	}

	utils.GinSuccessResponse(c, "Order retrieved successfully", order)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid user ID", err)
		return
	}

	orderID := c.Param("id")
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid order ID", err)
		return
	}

	err = h.orderService.CancelOrder(c.Request.Context(), orderUUID, userUUID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to cancel order", err)
		return
	}

	utils.GinSuccessResponse(c, "Order cancelled successfully", nil)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	utils.GinSuccessResponse(c, "All orders retrieved", []interface{}{})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	orderID := c.Param("id")
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid order ID", err)
		return
	}

	err = h.orderService.UpdateOrderStatus(c.Request.Context(), orderUUID, req.Status)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to update order status", err)
		return
	}

	utils.GinSuccessResponse(c, "Order status updated", nil)
}
