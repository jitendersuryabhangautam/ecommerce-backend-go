package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = b
	return rw.ResponseWriter.Write(b)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID
		requestID := uuid.New().String()

		// Wrap response writer to capture status code and body size
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Add request ID to header
		w.Header().Set("X-Request-ID", requestID)

		// Log request
		log.Printf("[%s] %s %s %s",
			requestID,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		// Process request
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log response
		log.Printf("[%s] %d %s %v",
			requestID,
			rw.statusCode,
			http.StatusText(rw.statusCode),
			duration,
		)
	})
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}
