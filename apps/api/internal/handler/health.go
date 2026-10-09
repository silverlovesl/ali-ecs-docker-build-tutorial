// Package handler contains the HTTP handlers exposed by the API.
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/config"
)

// HealthHandler exposes service liveness endpoints.
type HealthHandler struct {
	cfg config.Config
}

// NewHealthHandler builds a HealthHandler from the resolved configuration.
func NewHealthHandler(cfg config.Config) *HealthHandler {
	return &HealthHandler{cfg: cfg}
}

// Get reports service liveness and basic runtime metadata.
func (h *HealthHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "api",
		"env":     h.cfg.Env,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

// Hello is a sample endpoint consumed by the web app.
func (h *HealthHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello " + c.DefaultQuery("name", "world"),
	})
}
