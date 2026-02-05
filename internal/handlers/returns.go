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

type ReturnHandler struct {
	returnService service.ReturnService
}

func NewReturnHandler(returnService service.ReturnService) *ReturnHandler {
	return &ReturnHandler{returnService: returnService}
}

func (h *ReturnHandler) CreateReturn(c *gin.Context) {
	userID, err := middleware.GetUserIDFromGin(c)
	if err != nil {
		utils.GinUnauthorizedResponse(c, err.Error())
		return
	}

	var req models.CreateReturnRequest
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

	returnReq, err := h.returnService.CreateReturn(c.Request.Context(), req, userUUID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to create return", err)
		return
	}

	utils.GinCreatedResponse(c, "Return created successfully", returnReq)
}

func (h *ReturnHandler) GetUserReturns(c *gin.Context) {
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

	returns, total, err := h.returnService.GetUserReturns(c.Request.Context(), userUUID, page, limit)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to retrieve returns", err)
		return
	}

	response := map[string]interface{}{
		"returns": returns,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	}

	utils.GinSuccessResponse(c, "Returns retrieved successfully", response)
}

func (h *ReturnHandler) GetReturn(c *gin.Context) {
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

	returnID := c.Param("id")
	returnUUID, err := uuid.Parse(returnID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid return ID", err)
		return
	}

	returnReq, err := h.returnService.GetReturn(c.Request.Context(), returnUUID, userUUID)
	if err != nil {
		utils.GinNotFoundResponse(c, "Return not found")
		return
	}

	utils.GinSuccessResponse(c, "Return retrieved successfully", returnReq)
}

func (h *ReturnHandler) GetAllReturns(c *gin.Context) {
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

	status := c.Query("status")

	rangeDays := 0
	if rd := c.Query("range_days"); rd != "" {
		if parsed, err := strconv.Atoi(rd); err == nil && parsed > 0 {
			rangeDays = parsed
		}
	}

	returns, total, err := h.returnService.GetAllReturns(c.Request.Context(), page, limit, status, rangeDays)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to retrieve returns", err)
		return
	}

	response := map[string]interface{}{
		"returns": returns,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	}

	utils.GinSuccessResponse(c, "All returns retrieved", response)
}

func (h *ReturnHandler) ProcessReturn(c *gin.Context) {
	returnID := c.Param("id")
	if returnID == "" {
		returnID = c.Param("returnId")
	}
	returnUUID, err := uuid.Parse(returnID)
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid return ID", err)
		return
	}

	var req models.ProcessReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.GinBadRequestResponse(c, "Invalid request body", err)
		return
	}

	if errors := utils.ValidateStruct(req); errors != nil {
		utils.GinValidationErrorResponse(c, errors)
		return
	}

	returnReq, err := h.returnService.ProcessReturn(c.Request.Context(), returnUUID, req)
	if err != nil {
		utils.GinBadRequestResponse(c, "Failed to process return", err)
		return
	}

	utils.GinSuccessResponse(c, "Return processed", returnReq)
}
