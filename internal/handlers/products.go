package handlers

import (
	"net/http"
	"strconv"

	"ecommerce-backend/internal/models"
	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.ProductRequest

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

	// Create product
	product, err := h.productService.CreateProduct(c.Request.Context(), req)
	if err != nil {
		utils.GinErrorResponse(c, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	utils.GinCreatedResponse(c, "Product created successfully", product)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Get product ID from URL
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid product ID", err)
		return
	}

	// Get product
	product, err := h.productService.GetProduct(c.Request.Context(), productID)
	if err != nil {
		utils.GinNotFoundResponse(c, "Product")
		return
	}

	utils.GinSuccessResponse(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Get pagination parameters from query
	page := 1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	category := c.Query("category")
	search := c.Query("search")

	// Get products
	products, total, err := h.productService.GetProducts(c.Request.Context(), page, limit, category, search)
	if err != nil {
		utils.GinInternalErrorResponse(c, "Failed to get products", err)
		return
	}

	response := map[string]interface{}{
		"products": products,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	}

	utils.GinSuccessResponse(c, "Products retrieved successfully", response)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Get product ID from URL
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid product ID", err)
		return
	}

	var req models.ProductUpdateRequest

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

	// Update product
	product, err := h.productService.UpdateProduct(c.Request.Context(), productID, req)
	if err != nil {
		utils.GinErrorResponse(c, http.StatusBadRequest, "Failed to update product", err)
		return
	}

	utils.GinSuccessResponse(c, "Product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// Get product ID from URL
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.GinBadRequestResponse(c, "Invalid product ID", err)
		return
	}

	// Delete product
	err = h.productService.DeleteProduct(c.Request.Context(), productID)
	if err != nil {
		utils.GinErrorResponse(c, http.StatusBadRequest, "Failed to delete product", err)
		return
	}

	utils.GinSuccessResponse(c, "Product deleted successfully", nil)
}

func (h *ProductHandler) GetAdminProducts(c *gin.Context) {
	page := 1
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}

	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	rangeDays := 30
	if rd := c.Query("range_days"); rd != "" {
		if val, err := strconv.Atoi(rd); err == nil && val > 0 {
			rangeDays = val
		}
	}

	products, total, err := h.productService.GetAdminProducts(c.Request.Context(), page, limit, rangeDays)
	if err != nil {
		utils.GinInternalErrorResponse(c, "Failed to get products", err)
		return
	}

	response := map[string]interface{}{
		"products": products,
		"meta": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + limit - 1) / limit,
		},
	}

	utils.GinSuccessResponse(c, "Products retrieved successfully", response)
}

func (h *ProductHandler) GetTopProducts(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 100 {
			limit = val
		}
	}

	rangeDays := 0
	if rd := c.Query("range_days"); rd != "" {
		if val, err := strconv.Atoi(rd); err == nil && val > 0 {
			rangeDays = val
		}
	}

	items, err := h.productService.GetTopProducts(c.Request.Context(), limit, rangeDays)
	if err != nil {
		utils.GinInternalErrorResponse(c, "Failed to get top products", err)
		return
	}

	response := models.TopProductsResponse{
		RangeDays: rangeDays,
		Items:     items,
	}

	utils.GinSuccessResponse(c, "Top products retrieved", response)
}
