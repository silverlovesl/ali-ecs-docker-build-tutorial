// Package router wires middleware and routes into a Gin engine.
package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/config"
	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/handler"
	"github.com/example/ali-ecs-docker-build-tutorial/apps/api/internal/middleware"
)

// New builds a Gin engine with the project middleware and routes registered.
func New(cfg config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.AllowOrigins))

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	health := handler.NewHealthHandler(cfg)

	api := engine.Group("/api")
	{
		api.GET("/health", health.Get)
		api.GET("/hello", health.Hello)
	}

	return engine
}
