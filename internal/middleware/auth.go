package middleware

import (
	"context"
	"net/http"
	"strings"

	"ecommerce-backend/internal/service"
	"ecommerce-backend/pkg/utils"

	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	UserRoleKey contextKey = "userRole"
)

func AuthMiddleware(authService service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.UnauthorizedResponse(w)
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.UnauthorizedResponse(w)
				return
			}

			token := parts[1]

			// Validate token
			user, err := authService.ValidateToken(token)
			if err != nil {
				utils.UnauthorizedResponse(w)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, user.ID)
			ctx = context.WithValue(ctx, UserRoleKey, user.Role)

			// Continue with the request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user role from context
		role, ok := r.Context().Value(UserRoleKey).(string)
		if !ok || role != "admin" {
			utils.ForbiddenResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}
