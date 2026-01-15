package middleware

import (
	"log"
	"net/http"

	"ecommerce-backend/pkg/utils"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.Printf("PANIC: %v", err)

				// Return 500 Internal Server Error
				utils.ErrorResponse(w, http.StatusInternalServerError,
					"Internal server error",
					http.ErrAbortHandler)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
