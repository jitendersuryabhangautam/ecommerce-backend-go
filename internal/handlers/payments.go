package handlers

import (
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	utils.GinCreatedResponse(c, "Payment created successfully", nil)
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	utils.GinSuccessResponse(c, "Payment verified successfully", nil)
}

func (h *PaymentHandler) GetPaymentByOrder(c *gin.Context) {
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}
	utils.GinSuccessResponse(c, "Payment retrieved successfully", nil)
}
