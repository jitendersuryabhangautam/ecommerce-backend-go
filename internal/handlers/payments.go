package handlers

import (
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req models.CreatePaymentRequest
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

	payment, err := h.paymentService.CreatePayment(c.Request.Context(), req, userUUID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to create payment", err)
		return
	}

	utils.GinCreatedResponse(c, "Payment created successfully", payment)
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req models.VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	payment, err := h.paymentService.VerifyPayment(c.Request.Context(), req)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to verify payment", err)
		return
	}

	utils.GinSuccessResponse(c, "Payment verified successfully", payment)
}

func (h *PaymentHandler) GetPaymentByOrder(c *gin.Context) {
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	orderID := c.Param("id")
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid order ID", err)
		return
	}

	payment, err := h.paymentService.GetPaymentByOrderID(c.Request.Context(), orderUUID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to retrieve payment", err)
		return
	}

	if payment == nil {
		utils.GinNotFoundResponse(c, "Payment not found")
		return
	}

	utils.GinSuccessResponse(c, "Payment retrieved successfully", payment)
}
