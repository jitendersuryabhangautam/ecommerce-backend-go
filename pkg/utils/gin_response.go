package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinResponse is the standard response structure for Gin
type GinResponseData struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// GinSuccessResponse sends a success response with 200 status code
func GinSuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, GinResponseData{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// GinCreatedResponse sends a success response with 201 status code
func GinCreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, GinResponseData{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// GinErrorResponse sends an error response
func GinErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	c.JSON(statusCode, GinResponseData{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// GinValidationErrorResponse sends validation errors
func GinValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, GinResponseData{
		Success: false,
		Message: "Validation failed",
		Error:   "validation_error",
		Errors:  errors,
	})
}

// GinBadRequestResponse sends a 400 bad request response
func GinBadRequestResponse(c *gin.Context, message string, err error) {
	GinErrorResponse(c, http.StatusBadRequest, message, err)
}

// GinUnauthorizedResponse sends a 401 unauthorized response
func GinUnauthorizedResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	c.JSON(http.StatusUnauthorized, GinResponseData{
		Success: false,
		Message: message,
		Error:   "unauthorized",
	})
}

// GinForbiddenResponse sends a 403 forbidden response
func GinForbiddenResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	c.JSON(http.StatusForbidden, GinResponseData{
		Success: false,
		Message: message,
		Error:   "forbidden",
	})
}

// GinNotFoundResponse sends a 404 not found response
func GinNotFoundResponse(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, GinResponseData{
		Success: false,
		Message: resource + " not found",
		Error:   "not_found",
	})
}

// GinConflictResponse sends a 409 conflict response
func GinConflictResponse(c *gin.Context, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	c.JSON(http.StatusConflict, GinResponseData{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}

// GinInternalErrorResponse sends a 500 internal server error response
func GinInternalErrorResponse(c *gin.Context, message string, err error) {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	c.JSON(http.StatusInternalServerError, GinResponseData{
		Success: false,
		Message: message,
		Error:   errorMsg,
	})
}
