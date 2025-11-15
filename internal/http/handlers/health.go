package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
	pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	// Check database connection
	err := h.pool.Ping(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "error",
			"database": "error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "ok",
	})
}




