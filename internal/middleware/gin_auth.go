package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"ecommerce-backend/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	GinUserIDKey   = "userID"
	GinUserRoleKey = "userRole"
)

// GinAuthMiddleware validates JWT tokens for Gin
func GinAuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Missing authorization header",
			})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Invalid token format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		user, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   "Invalid token",
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set(GinUserIDKey, user.ID.String())
		c.Set(GinUserRoleKey, user.Role)

		c.Next()
	}
}

// GinAdminMiddleware checks if user has admin role
func GinAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(GinUserRoleKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "User role not found",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden",
				"error":   "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserIDFromGin extracts user ID from Gin context
func GetUserIDFromGin(c *gin.Context) (string, error) {
	userID, exists := c.Get(GinUserIDKey)
	if !exists {
		return "", ErrUserIDNotFound
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("invalid user ID type")
	}
	return userIDStr, nil
}

// GetUserRoleFromGin extracts user role from Gin context
func GetUserRoleFromGin(c *gin.Context) (string, error) {
	role, exists := c.Get(GinUserRoleKey)
	if !exists {
		return "", ErrUserRoleNotFound
	}

	return role.(string), nil
}

var (
	ErrUserIDNotFound   = NewAPIError("user_id_not_found", "User ID not found in context")
	ErrUserRoleNotFound = NewAPIError("user_role_not_found", "User role not found in context")
)

// NewAPIError creates an error with code and message
func NewAPIError(code, message string) error {
	return &apiError{code: code, message: message}
}

type apiError struct {
	code    string
	message string
}

func (e *apiError) Error() string {
	return e.message
}
