package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	db *pgxpool.Pool
}

func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "ecommerce-backend",
		"version":   "1.0.0",
	}

	// Check database connection
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		health["status"] = "unhealthy"
		health["database"] = "disconnected"
		health["error"] = err.Error()
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}

	health["database"] = "connected"
	c.JSON(http.StatusOK, health)
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	readiness := map[string]interface{}{
		"ready":     true,
		"timestamp": time.Now().UTC(),
	}

	// Check database
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		readiness["ready"] = false
		readiness["database"] = "not_ready"
		readiness["error"] = err.Error()
		c.JSON(http.StatusServiceUnavailable, readiness)
		return
	}

	readiness["database"] = "ready"
	c.JSON(http.StatusOK, readiness)
}

func (h *HealthHandler) Metrics(c *gin.Context) {
	metrics := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"metrics": map[string]interface{}{
			"uptime": time.Since(startTime).String(),
		},
	}

	c.JSON(http.StatusOK, metrics)
}

var startTime = time.Now()
