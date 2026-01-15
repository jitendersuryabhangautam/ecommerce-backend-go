package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func SuccessResponse(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSONResponse(w, http.StatusOK, response)
}

func CreatedResponse(w http.ResponseWriter, message string, data interface{}) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	JSONResponse(w, http.StatusCreated, response)
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string, err error) {
	response := Response{
		Success: false,
		Message: message,
		Error:   err.Error(),
	}
	JSONResponse(w, statusCode, response)
}

func ValidationErrorResponse(w http.ResponseWriter, errors map[string]string) {
	response := Response{
		Success: false,
		Message: "Validation failed",
		Error:   "validation_error",
		Data:    errors,
	}
	JSONResponse(w, http.StatusBadRequest, response)
}

func NotFoundResponse(w http.ResponseWriter, resource string) {
	response := Response{
		Success: false,
		Message: resource + " not found",
		Error:   "not_found",
	}
	JSONResponse(w, http.StatusNotFound, response)
}

func UnauthorizedResponse(w http.ResponseWriter) {
	response := Response{
		Success: false,
		Message: "Unauthorized",
		Error:   "unauthorized",
	}
	JSONResponse(w, http.StatusUnauthorized, response)
}

func ForbiddenResponse(w http.ResponseWriter) {
	response := Response{
		Success: false,
		Message: "Forbidden",
		Error:   "forbidden",
	}
	JSONResponse(w, http.StatusForbidden, response)
}
