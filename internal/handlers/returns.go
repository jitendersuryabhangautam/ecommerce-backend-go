package handlers

import (
	"ecommerce-backend/internal/middleware"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ReturnHandler struct {
	returnService service.ReturnService
}

func NewReturnHandler(returnService service.ReturnService) *ReturnHandler {
	return &ReturnHandler{returnService: returnService}
}

func (h *ReturnHandler) CreateReturn(c *gin.Context) {
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

	utils.GinCreatedResponse(c, "Return created successfully", nil)
}

func (h *ReturnHandler) GetUserReturns(c *gin.Context) {
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}
	utils.GinSuccessResponse(c, "Returns retrieved successfully", []interface{}{})
}

func (h *ReturnHandler) GetReturn(c *gin.Context) {
	_, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}
	utils.GinSuccessResponse(c, "Return retrieved successfully", nil)
}

func (h *ReturnHandler) GetAllReturns(c *gin.Context) {
	utils.GinSuccessResponse(c, "All returns retrieved", []interface{}{})
}

func (h *ReturnHandler) ProcessReturn(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}
	utils.GinSuccessResponse(c, "Return processed", nil)
}
