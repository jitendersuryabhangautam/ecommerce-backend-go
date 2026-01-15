package middleware

import (
	"net/http"
	"sync"
	"time"

	"ecommerce-backend/pkg/utils"
)

type rateLimiter struct {
	visits map[string][]time.Time
	mu     sync.RWMutex
}

func newRateLimiter() *rateLimiter {
	return &rateLimiter{
		visits: make(map[string][]time.Time),
	}
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, visits := range rl.visits {
		var validVisits []time.Time
		for _, visit := range visits {
			if time.Since(visit) < time.Minute {
				validVisits = append(validVisits, visit)
			}
		}
		if len(validVisits) == 0 {
			delete(rl.visits, ip)
		} else {
			rl.visits[ip] = validVisits
		}
	}
}

func (rl *rateLimiter) allow(ip string, limit int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	visits := rl.visits[ip]

	// Remove old visits
	var recentVisits []time.Time
	for _, visit := range visits {
		if now.Sub(visit) < time.Minute {
			recentVisits = append(recentVisits, visit)
		}
	}

	// Check if limit exceeded
	if len(recentVisits) >= limit {
		return false
	}

	// Add new visit
	recentVisits = append(recentVisits, now)
	rl.visits[ip] = recentVisits

	return true
}

func RateLimit(limit int) func(http.Handler) http.Handler {
	limiter := newRateLimiter()

	// Clean up old entries periodically
	go func() {
		for {
			time.Sleep(time.Minute)
			limiter.cleanup()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			// Check rate limit
			if !limiter.allow(ip, limit) {
				utils.ErrorResponse(w, http.StatusTooManyRequests,
					"Rate limit exceeded",
					http.ErrAbortHandler)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Per-user rate limiting
func UserRateLimit(limit int) func(http.Handler) http.Handler {
	limiter := newRateLimiter()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get user ID from context for authenticated users
			userID, ok := GetUserIDFromContext(r.Context())
			if ok {
				// Use user ID for rate limiting
				if !limiter.allow(userID.String(), limit) {
					utils.ErrorResponse(w, http.StatusTooManyRequests,
						"Rate limit exceeded",
						http.ErrAbortHandler)
					return
				}
			} else {
				// Use IP for unauthenticated users
				ip := r.RemoteAddr
				if !limiter.allow(ip, limit) {
					utils.ErrorResponse(w, http.StatusTooManyRequests,
						"Rate limit exceeded",
						http.ErrAbortHandler)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
