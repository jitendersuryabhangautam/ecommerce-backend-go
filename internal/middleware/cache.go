package middleware

import (
	"github.com/gin-gonic/gin"
)

// NoCacheMiddleware sets cache control headers to prevent caching
func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent caching of API responses
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}
}
